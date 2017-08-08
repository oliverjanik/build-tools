package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: hc <html-file> {extract|inject} [-i input] [-o output]")
	}

	path := os.Args[1]
	cmd := os.Args[2]
	rest := os.Args[3:]

	switch cmd {
	case "extract":
		flags := flag.NewFlagSet(cmd, flag.ExitOnError)
		output := flags.String("o", "", "Output file")
		flags.Parse(rest)

		extract(path, *output)
	case "inject":
	default:
		log.Fatalln("Unknown command", cmd)
	}
}

func extract(path string, output string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	root, err := html.Parse(f)
	if err != nil {
		log.Fatalln("Error parsing", path, ":", err)
	}

	confNode := findConfig(root)
	if confNode == nil || confNode.FirstChild == nil || confNode.FirstChild.Type != html.TextNode {
		return
	}

	result := confNode.FirstChild.Data

	if output == "" {
		fmt.Print(result)
		return
	}

	fo, err := os.Create(output)
	if err != nil {
		log.Fatalln(err)
	}

	defer fo.Close()

	_, err = fo.WriteString(result)
	if err != nil {
		log.Fatalln(err)
	}
}

func findConfig(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	if n.Data == "script" && getID(n) == "config" {
		return n
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if found := findConfig(child); found != nil {
			return found
		}
	}

	return nil
}

func getID(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}

	return ""
}
