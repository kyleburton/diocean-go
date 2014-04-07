/*
Package diocean provides a set of types and functions to interact with Digital Ocean's API: https://developers.digitalocean.com/

Configuration

Create the file ~/.digitalocean.json with the following content:

    {
      "ClientId": "<<your client-id>>",
      "ApiKey":   "<<your api-key>>"
    }

You can find your ClientId and ApiKey on the api_access page for your account https://cloud.digitalocean.com/api_access

Basics

Create a client struct with your ClientId and ApiKey:

    Client = &diocean.DioceanClient{
      ClientId:      "client-id",
      ApiKey:        "api-key",
      Verbose:       false,
      WaitForEvents: true,
    }

You may optionally supply a Verbose flag and WaitForEvents flag.


*/
package diocean

/*

*/
