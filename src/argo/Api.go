package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func getToken(username string, password string, host string) string {

	client := buildHttpClient()

	message := map[string]interface{}{
		"username": username,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Post(host+"/api/v1/session", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result["token"].(string)
}

func GetResourceTree(applicationName string, host string) {
	token := getToken("admin", "newpassword", "https://34.71.103.174")

	client := buildHttpClient()

	req, err := http.NewRequest("GET", "https://34.71.103.174/api/v1/applications/"+applicationName+"/resource-tree", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	result, _ := ioutil.ReadAll(resp.Body)

	log.Print(string(result))
}

//async getManagedResources() {
//const resourceTree = await rp({
//method: 'GET',
//uri: `${host}/api/v1/applications/${applicationName}/resource-tree`,
//headers: {
//'Authorization': `Bearer ${this.token}`
//},
//json: true
//});
//
//const deploymentStatuses = resourceTree.nodes.filter(node => node.kind === 'Deployment').map(deploy => {
//return {
//id: deploy.uid,
//status: _.get(deploy.health, 'status', 'Missing'),
//}
//});
//
//const result = await rp({
//method: 'GET',
//uri: `${host}/api/v1/applications/${applicationName}/managed-resources`,
//headers: {
//'Authorization': `Bearer ${this.token}`
//},
//json: true
//});
//
//const deployments = _.filter(result.items, (item) => item.kind === 'Deployment');
//return deployments.map((deployment) => {
//const targetState = JSON.parse(deployment.targetState);
//const liveState = JSON.parse(deployment.liveState);
//
//const result = {
//name: deployment.name
//};
//
//if(targetState) {
//result.targetImages = targetState.spec.template.spec.containers.map(container => container.image);
//}
//
//if(liveState) {
//result.status = deploymentStatuses.find(deploymentStatus => deploymentStatus.id === liveState.metadata.uid).status;
//result.liveImages = liveState.spec.template.spec.containers.map(container => container.image);
//}
//return result;
//
//});
//}
