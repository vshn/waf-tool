package forwarder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const loggingNamespace = "openshift-logging"

// PortForwarder can forward a port to Elasticsearch
type PortForwarder interface {
	Start() error
	Stop() error
}

type portForwarder struct {
	cmd  *exec.Cmd
	port string
}

// New creates a new port forwarder
func New(port string) PortForwarder {
	return &portForwarder{
		port: port,
	}
}

// Start starts a port forwarding
func (p *portForwarder) Start() error {

	if p.cmd != nil {
		return errors.New("this port forwarder is already in use")
	}
	p.cmd = exec.Command("oc", "port-forward", "--namespace", loggingNamespace, "svc/logging-es", p.port)
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	errorChan := make(chan error)

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			errorChan <- errors.New(line)
			return
		}
		errorChan <- scanner.Err()
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Debug(line)
			errorChan <- nil
		}
		errorChan <- nil
	}()

	log.WithField("port", p.port).Info("Starting port forward...")
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("could not start port forwarding: %w", err)
	}

	if err := <-errorChan; err != nil {
		return err
	}

	return nil
}

// Stop the forwarding
func (p *portForwarder) Stop() error {
	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	err := p.cmd.Wait()
	if err != nil && !strings.HasPrefix(err.Error(), "exit status") {
		return err
	}
	return nil
}
