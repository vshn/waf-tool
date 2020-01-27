# waf-tuning-tool
Helper tool to generate rule exclusions based on ElasticSearch logs queried by uniqueID of the access log

## Concept
The first version of the tool would have to operate as follows:

1. Authenticate in OpenShift via a service account using either a token or acc/pwd credentials

    <code>oc login</code>

2. Since there are no routes to access the Elasticsearch directly a port forwarding is necessary. Create a port forwarding from the local client to the service of Elasticsearch in the Openshift cluster. Elasticsearch is listening on port 9200.

    <code>oc port-forward -n logging svc/logging-es 9200:9200</code>

3. Access the logs of elasticsearch and search for the desired modsec-alert uniqueId that was given as the input of the software. All the indexes that begin with project.* have to be searched.

    <pre><code>curl -XGET "https://logging-es:9200/project.*/_search" -H 'Content-Type: application/json' -d'
    {
      "_source": ["modsec-alert"]
      "query": {
        "bool": {
          "must": [
            {
              "match_phrase": {
                "modsec-alert.uniqueID": {
                  "query": $uniqueId
                }
              }
            }
          ]
        }
      }
    }'</code></pre>

4. The output would be a json that has to be parsed, recommended by a library. Example output.

    <pre><code>"modsec-alert": {
      "msg": "Restricted SQL Character Anomaly Detection (args): # of special characters exceeded (12)",
      "severity": "WARNING",
      "ver": "OWASP_CRS/3.1.0",
      the comment in the conf file
      "rule_template": "# ModSec Rule Exclusion: 942430 : Restricted SQL Character Anomaly Detection (args): # of special characters exceeded
      (12) (severity:  WARNING)",
      "data": "Matched Data: -0620-41bb-8c2f-d024922c173e.6404161c-05af-498f-a222-eceece5bf4ab.0ce081c4-3353-4b18-a764- found within ARGS:code: c3946ac9-0620-41bb-8c2f-d024922c173e.6404161c-05af-498f-a222-eceece5bf4ab.0ce081c4-3353-4b18-a764-8a95631a6e9c",
      "line": 1235,
      "description": """ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98`<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98`<>]*?){12})" at ARGS:code.""",
      "uri": "/oauth2-redirect.html",
      "tags": [
        "application-multi",
        "language-multi",
        "platform-multi",
        "attack-sqli",
        "OWASP_CRS/WEB_ATTACK/SQL_INJECTION",
        "WASCTC/WASC-19",
        "OWASP_TOP_10/A1",
        "OWASP_AppSensor/CIE1",
        "PCI/6.5.2",
        "paranoia-level/2"
      ],
      "hostname": "localhost",
      "file": "/etc/apache2/modsecurity.d/owasp-crs/rules/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
      "client": "85.195.222.81",
      "id": 942430,
      "uniqueID": "Xibh0U-m3UNWcN42z4ArFQAAAFA"
    }</code></pre>

5. Parse the relevant information from the modsec-alert object like id and ARGS variable from description.
6. Generate several exclusion rules for the end user to choose from, examples:
    <pre><code>1. SecRule REQUEST_URI "@beginsWith /auth/oidc/callback" "phase:2,nolog,pass,ctl:ruleRemoveTargetById=942430;ARGS:code" (most frequently used)
   1. SecRule REQUEST_URI "@beginsWith /auth/oidc/callback" "phase:2,nolog,pass,ctl:ruleRemoveById=942430"</code></pre>

## More Info
For more details access the wiki page [WAF Tuning Tool](https://wiki.vshn.net/display/FILO/WAF+Tuning+Tool)
