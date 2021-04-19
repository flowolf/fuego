package main

import (
	firestore "cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/urfave/cli"
)

func addData(
	client *firestore.Client,
	collection string,
	data string,
	timestampify bool,
	name string) (string, error) {

	object, err := unmarshallData(data)
	if err != nil {
		return "", err
	}

	if timestampify {
		timestampifyMap(object)
	}

	if name != "" {
		_, err := client.
		Collection(collection).Doc(name).
		Set(context.Background(), object);
		if err != nil {
			return "", err
		}
		return name, nil
	} else {
		doc, _, err := client.
			Collection(collection).
			Add(context.Background(), object)
		if err != nil {
			return "", err
		}
		return doc.ID, nil
	}
}

func addCommandAction(c *cli.Context) error {
	collectionPath := c.Args().First()
	timestampify := c.Bool("timestamp")
	name := c.String("name")
	data := c.Args().Get(1)

	client, err := createClient(credentials)
	if err != nil {
		return cliClientError(err)
	}
	id, err := addData(client, collectionPath, data, timestampify, name)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to add data. \n%v", err), 81)
	}
	fmt.Fprintf(c.App.Writer, "%v\n", id)
	defer client.Close()
	return nil
}
