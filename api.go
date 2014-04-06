package diocean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type DioceanClient struct {
  ClientId string
  ApiKey string
  Verbose bool
  WaitForEvents bool
}

// https://developers.digitalocean.com/sizes/

// GET https://api.digitalocean.com/sizes/?client_id=[client_id]&api_key=[api_key]

////////////////////

type DropletSizesResponse struct {
	Status string
	Sizes  []DropletSize
}

func (self *DropletSizesResponse) Header() []string {
	return []string{
		"id",
		"name",
		"slug",
	}
}

func (self *DropletSizesResponse) Unmarshal(content []byte) {
	json.Unmarshal(content, self)
	if self.Sizes == nil {
		self.Sizes = make([]DropletSize, 0)
	}
}

////////////////////

type DropletSize struct {
	Id   float64
	Name string
	Slug string
}

func (self *DropletSize) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
		self.Slug,
	}
}

////////////////////

type ActiveDropletsResponse struct {
	Status   string
	Droplets []DropletInfo
}

func (self *ActiveDropletsResponse) Header() []string {
	return []string{
		"id",
		"name",
		"image_id",
		"size_id",
		"region_id",
		"backups_active",
		"ip_address",
		"private_ip_address",
		"locked",
		"status",
		"created_at",
	}
}

func (self *ActiveDropletsResponse) Unmarshal(content []byte) {
	json.Unmarshal(content, self)

	if self.Droplets == nil {
		self.Droplets = make([]DropletInfo, 0)
	}
}

type DropletShowResponse struct {
	Status  string
	Droplet ShowDropletInfo
}

func (self *DropletShowResponse) Header() []string {
	return []string{
		"id",
		"image_id",
		"name",
		"region_id",
		"size_id",
		"backups_active",
		"backups",
		"snapshots",
		"ip_address",
		"private_ip_address",
		"locked",
		"status",
	}
}

func (self *DropletShowResponse) Unmarshal(content []byte) {
	json.Unmarshal(content, self)

	// if self.Droplet.backups == nil {
	//   self.Droplet.backups = make([]???)
	// }

	// if self.Droplet.snapshots == nil {
	//   self.Droplet.snapshots = make([]???)
	// }
}

type ShowDropletInfo struct {
	Id                 float64
	Image_id           float64
	Name               string
	Region_id          float64
	Size_id            float64
	Backups_active     bool
	Backups            []interface{}
	Snapshots          []interface{}
	Ip_address         string
	Private_ip_address string
	Locked             bool
	Status             string
}

func (self *ShowDropletInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		fmt.Sprintf("%.f", self.Image_id),
		self.Name,
		fmt.Sprintf("%.f", self.Region_id),
		fmt.Sprintf("%.f", self.Size_id),
		fmt.Sprintf("%t", self.Backups_active),
		fmt.Sprintf("%d", self.Backups),
		fmt.Sprintf("%d", self.Snapshots),
		self.Ip_address,
		self.Private_ip_address,
		fmt.Sprintf("%t", self.Locked),
		self.Status,
	}
}

type DropletInfo struct {
	Id                 float64
	Name               string
	Image_id           float64
	Size_id            float64
	Region_id          float64
	Backups_active     bool
	Ip_address         string
	Private_ip_address string
	Locked             bool
	Status             string
	Created_at         string
}

func (self *DropletInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
		fmt.Sprintf("%.f", self.Image_id),
		fmt.Sprintf("%.f", self.Size_id),
		fmt.Sprintf("%.f", self.Region_id),
		fmt.Sprintf("%t", self.Backups_active),
		self.Ip_address,
		self.Private_ip_address,
		fmt.Sprintf("%t", self.Locked),
		self.Status,
		self.Created_at,
	}
}

////////////////////

type NewDropletResponse struct {
	Status  string
	Droplet NewDropletInfo
}

func (self *NewDropletResponse) Header() []string {
	return []string{
		"id",
		"name",
		"image_id",
		"size_id",
		"event_id",
	}
}

func (self *NewDropletResponse) Unmarshal(content []byte) {
	json.Unmarshal(content, self)
}

type NewDropletInfo struct {
	Id       float64
	Name     string
	Image_id float64
	Size_id  float64
	Event_id float64
}

func (self *NewDropletInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
		fmt.Sprintf("%.f", self.Image_id),
		fmt.Sprintf("%.f", self.Size_id),
		fmt.Sprintf("%.f", self.Event_id),
	}
}

////////////////////

type SimpleResponse struct {
	Status string
}

func (self *SimpleResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
}

type SimpleEventResponse struct {
	Status   string
	Event_id float64
}

func (self *SimpleEventResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
}

func (self *SimpleEventResponse) Header() []string {
	return []string{
		"event_id",
	}
}

////////////////////

type RegionResponse struct {
	Status  string
	Regions []RegionInfo
}

func (self *RegionResponse) Header() []string {
	return []string{
		"id",
		"name",
		"slug",
	}
}

func (self *RegionResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)

	if self.Regions == nil {
		self.Regions = make([]RegionInfo, 0)
	}
}

type RegionInfo struct {
	Id   float64
	Name string
	Slug string
}

func (self *RegionInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
		self.Slug,
	}
}

////////////////////

type SshKeysResponse struct {
	Status   string
	Ssh_keys *[]SshKeyInfo
}

func (self *SshKeysResponse) Header() []string {
	return []string{
		"Id",
		"Name",
	}
}

func (self *SshKeysResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
	if self.Ssh_keys == nil {
		keys := make([]SshKeyInfo, 0)
		self.Ssh_keys = &keys
	}
}

type SshKeyInfo struct {
	Id   float64
	Name string
}

func (self *SshKeyInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
	}
}

////////////////////

type ImagesResponse struct {
	Id     float64
	Images []ImageInfo
}

func (self *ImagesResponse) Header() []string {
	return []string{
		"id",
		"name",
		"distribution",
		"slug",
		"public",
	}
}

func (self *ImagesResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
}

type ImageShowResponse struct {
	Status string
	Image  ImageInfo
}

func (self *ImageShowResponse) Header() []string {
	return []string{
		"id",
		"name",
		"distribution",
		"slug",
		"public",
	}
}

func (self *ImageShowResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
}

type ImageInfo struct {
	Id           float64
	Name         string
	Distribution string
	Slug         string
	Public       bool
}

func (self *ImageInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Name,
		self.Distribution,
		self.Slug,
		fmt.Sprintf("%t", self.Public),
	}
}

////////////////////

type EventResponse struct {
	Status string
	Event  EventInfo
}

func (self *EventResponse) Unmarshal(body []byte) {
	json.Unmarshal(body, self)
}

func (self *EventResponse) Header() []string {
	return []string{
		"id",
		"action_status",
		"droplet_id",
		"event_type_id",
		"percentage",
	}
}

type EventInfo struct {
	Id            float64
	Action_status string
	Droplet_id    float64
	Event_type_id float64
	Percentage    string
}

func (self *EventInfo) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%.f", self.Id),
		self.Action_status,
		fmt.Sprintf("%.f", self.Droplet_id),
		fmt.Sprintf("%.f", self.Event_type_id),
		self.Percentage,
	}
}

////////////////////

func (self *DioceanClient) ApiGet(path string, params *url.Values) (*http.Response, []byte, error) {
	if params == nil {
		params = &url.Values{}
	}
	params.Add("client_id", self.ClientId)
	params.Add("api_key", self.ApiKey)
	url := fmt.Sprintf("https://api.digitalocean.com%s?", path) + params.Encode()
	if self.Verbose {
		fmt.Fprintf(os.Stderr, "ApiGet: url=%s\n", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error[GET:%s]: %s\n", url, err)
		return resp, nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading body: %s\n", err)
		return resp, nil, err
	}

	return resp, body, nil
}

func MapGetString(m map[string]interface{}, k string, defaultValue string) string {
	val, ok := m[k]
	if ok && val != nil {
		return val.(string)
	}

	return defaultValue
}

func (self *DioceanClient) DoImagesLs() {
	path := "/images/"
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp ImagesResponse
	resp.Unmarshal(body)

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")

	for _, image := range resp.Images {
		fmt.Print(strings.Join(image.ToStringArray(), "\t"))
		fmt.Print("\n")
	}
}

func (self *DioceanClient) DoImageShow(image_id string) {
	path := fmt.Sprintf("/images/%s", image_id)
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp ImageShowResponse
	resp.Unmarshal(body)
	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	fmt.Print(strings.Join(resp.Image.ToStringArray(), "\t"))
	fmt.Print("\n")
}

func (self *DioceanClient) DoImageDestroy(image_id string) {
	path := fmt.Sprintf("/images/%s/destroy", image_id)
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp SimpleResponse
	resp.Unmarshal(body)

	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Printf("Image Destroyed: %s\n", resp.Status)

}

func (self *DioceanClient) EventShow(eventId string) *EventResponse {
	path := fmt.Sprintf("/events/%s/", eventId)
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	if self.Verbose {
		fmt.Fprintf(os.Stderr, "body=%s\n", body)
	}

	var resp EventResponse
	resp.Unmarshal(body)
	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	if self.Verbose {
		fmt.Fprintf(os.Stderr, "resp=%s\n", resp)
	}

	return &resp
}

func (self *DioceanClient) DoEventShow(event_id string) {
	resp := self.EventShow(event_id)

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	fmt.Print(strings.Join(resp.Event.ToStringArray(), "\t"))
	fmt.Print("\n")
}

func (self *DioceanClient) WaitForEvent(eventId string) {
	resp := self.EventShow(eventId)

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	for {
		fmt.Print(strings.Join(resp.Event.ToStringArray(), "\t"))
		fmt.Print("\n")
		if resp.Event.Percentage == "100" {
			break
		}
		resp = self.EventShow(eventId)
	}
}

func (self *DioceanClient) DoEventWait(event_id string) {
	self.WaitForEvent(event_id)
}

func (self *DioceanClient) DropletsLs() *ActiveDropletsResponse {
	path := "/droplets/"
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp ActiveDropletsResponse
	resp.Unmarshal(body)

	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	return &resp
}

func (self *DioceanClient) DoDropletsLs() {
	resp := self.DropletsLs()

	fmt.Printf("%s\n", strings.Join(resp.Header(), "\t"))
	for _, droplet := range resp.Droplets {
		fmt.Print(strings.Join(droplet.ToStringArray(), "\t"))
		fmt.Print("\n")
	}
}

func (self *DioceanClient) DoApiGetWithSimpleResponse(path string, params *url.Values) {
	_, body, err := self.ApiGet(path, params)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp SimpleEventResponse
	resp.Unmarshal(body)

	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	fmt.Printf("%.f", resp.Event_id)
	fmt.Print("\n")

	if self.WaitForEvents {
		self.WaitForEvent(fmt.Sprintf("%.f", resp.Event_id))
	}
}

func (self *DioceanClient) DoDropletsDestroyDroplet(droplet_id string, scrub_data string) {
	params := &url.Values{}
	if self.Verbose {
		fmt.Printf("DoDropletsDestroyDroplet: droplet_id=%s\n", droplet_id)
	}
	params.Add("scrub_data", scrub_data)

	path := fmt.Sprintf("/droplets/%s/destroy/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, params)
}

func ParamsAddSize(params *url.Values, size string) {
	matched, err := regexp.MatchString("^\\d+$", size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error regex match failed: %s\n", err)
		os.Exit(1)
	}

	if matched {
		params.Add("size_id", size)
	} else {
		params.Add("size_slug", size)
	}
}

func (self *DioceanClient) DoDropletsNewDroplet (name string, size string, image string, region string, ssh_key_ids string, private_networking string, backups_enabled string) {
	params := &url.Values{}
	params.Add("name", name)

	ParamsAddSize(params, size)

	matched, err := regexp.MatchString("^\\d+$", image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error regex match failed: %s\n", err)
		os.Exit(1)
	}

	if matched {
		params.Add("image_id", image)
	} else {
		params.Add("image_slug", image)
	}

	matched, err = regexp.MatchString("^\\d+$", region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error regex match failed: %s\n", err)
		os.Exit(1)
	}

	if matched {
		params.Add("region_id", region)
	} else {
		params.Add("region_slug", region)
	}

	params.Add("ssh_key_ids", ssh_key_ids)
	params.Add("private_networking", private_networking)
	params.Add("backups_enabled", backups_enabled)

	_, body, err := self.ApiGet("/droplets/new", params)

	var resp NewDropletResponse
	resp.Unmarshal(body)

	if resp.Status != "OK" {
    fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	fmt.Print(strings.Join(resp.Droplet.ToStringArray(), "\t"))
	fmt.Print("\n")

}

func (self *DioceanClient) DoDropletsLsDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Fprintf(os.Stderr, "DoDropletsLsDroplet %s\n", droplet_id)
	}
	path := fmt.Sprintf("/droplets/%s", droplet_id)
	_, body, err := self.ApiGet(path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp DropletShowResponse
	resp.Unmarshal(body)

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	fmt.Print(strings.Join(resp.Droplet.ToStringArray(), "\t"))
	fmt.Print("\n")
}

func (self *DioceanClient) DoDropletsRebootDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsRebootDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/reboot/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsPowerCycleDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsPowerCycleDroplet: droplet_id=%s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/power_cycle/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsShutDownDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsShutDownDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/shutdown/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsPowerOffDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsPowerOffDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/power_off/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsPasswordResetDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsPasswordResetDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/password_reset/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsResizeDroplet(droplet_id string, size string) {
	if self.Verbose {
		fmt.Printf("DoDropletsResizeDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/resize/", droplet_id)
	params := &url.Values{}
	ParamsAddSize(params, size)
	self.DoApiGetWithSimpleResponse(path, params)
}

func (self *DioceanClient) DoDropletsSnapshotDroplet(droplet_id string, name string) {
	if self.Verbose {
		fmt.Printf("DoDropletsSnapshotDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/snapshot/", droplet_id)
	params := &url.Values{}
  params.Add("name", name)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DoDropletsPowerOnDroplet(droplet_id string) {
	if self.Verbose {
		fmt.Printf("DoDropletsPowernfDroplet: %s\n", droplet_id)
	}

	path := fmt.Sprintf("/droplets/%s/power_on/", droplet_id)
	self.DoApiGetWithSimpleResponse(path, nil)
}

func (self *DioceanClient) DropletSizesLs () {
	path := "/sizes/"
	_, body, err := self.ApiGet(path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp DropletSizesResponse
	resp.Unmarshal(body)
	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	for _, size := range resp.Sizes {
		fmt.Print(strings.Join(size.ToStringArray(), "\t"))
		fmt.Print("\n")
	}
}

func (self *DioceanClient) DoRegionsLs() {

	path := "/regions/"
	_, body, err := self.ApiGet(path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	var resp RegionResponse
	resp.Unmarshal(body)
	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	for _, region := range resp.Regions {
		fmt.Print(strings.Join(region.ToStringArray(), "\t"))
		fmt.Print("\n")
	}
}

func (self *DioceanClient) DoSshKeysLs() {
	path := "/ssh_keys/"
	_, body, err := self.ApiGet(path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing http.get[%s]: %s\n", path, err)
		os.Exit(1)
	}

	if self.Verbose {
		fmt.Fprintf(os.Stderr, "body=%s\n", body)
	}

	var resp SshKeysResponse
	resp.Unmarshal(body)
	if resp.Status != "OK" {
		fmt.Fprintf(os.Stderr, "Error: status != OK status=%s resp=%s\n", resp.Status, string(body))
		os.Exit(1)
	}

	if self.Verbose {
		fmt.Fprintf(os.Stderr, "resp=%s\n", resp)
	}

	fmt.Print(strings.Join(resp.Header(), "\t"))
	fmt.Print("\n")
	for _, sshKey := range *resp.Ssh_keys {
		fmt.Print(strings.Join(sshKey.ToStringArray(), "\t"))
		fmt.Print("\n")
	}
}

func (self *DioceanClient) DoSshFixKnownHosts() {
	resp := self.DropletsLs()

	for _, droplet := range resp.Droplets {
		cmd := fmt.Sprintf("ssh-keygen -f %s/.ssh/known_hosts -R %s", os.Getenv("HOME"), droplet.Ip_address)
		fmt.Printf("%s\n", cmd)
		out, err := exec.Command(
			"ssh-keygen",
			"-f",
			fmt.Sprintf("%s/.ssh/known_hosts", os.Getenv("HOME")),
			"-R",
			droplet.Ip_address,
		).Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing cmd[%s] : %s\n", cmd, err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	}
}

