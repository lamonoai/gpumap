package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var filename string

	flag.StringVar(&filename, "f", "", "")
	flag.Parse()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	sl := make([]string, 0)
	for scanner.Scan() {
		sl = append(sl, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	nl, err := formatLines(sl)
	if err != nil {
		log.Fatal(err)
	}

	nf, err := outputLines(filename, nl)
	if err != nil {
		if len(nf) != 0 {
			now := time.Now().String()
			tmp := fmt.Sprintf("%s.%s.%s", filename, "tmp", now)
			if err := os.Rename(filename, tmp); err != nil {
				log.Fatal(err)
			}

			if err := os.Rename(nf, filename); err != nil {
				log.Fatal(err)
			}
		}
		log.Fatal(err)
	} else {
		if err := os.Remove(nf); err != nil {
			log.Fatal(err)
		}
	}
}

func formatLines(sl []string) ([]string, error) {
	os := make([]string, 0)
	for _, s := range sl {
		ns, err := formatLine(s)
		if err != nil {
			return nil, err
		}
		os = append(os, ns)
	}

	return os, nil
}

func formatLine(s string) (string, error) {
	s = strings.TrimSpace(s)
	if !strings.Contains(s, "|") {
		return s, nil
	}

	if s[0] != '|' || s[len(s)-1] != '|' {
		return "", fmt.Errorf("line format error %s", s)
	}

	l := make([]string, 0)
	l = append(l, "|")

	w := s[1:]
	for {
		if len(w) == 0 {
			break
		}

		i := strings.Index(w, "|")
		if i == -1 {
			break
		}

		t := w[:i]
		e := strings.TrimSpace(t)
		oe := fmt.Sprintf("%s%s%s%s", " ", e, " ", "|")
		l = append(l, oe)

		w = w[i+1:]
	}

	ns := strings.Join(l, "")

	return ns, nil
}

func outputLines(filename string, sl []string) (string, error) {
	now := time.Now().String()
	nf := fmt.Sprintf("%s.%s", filename, now)
	if err := os.Rename(filename, nf); err != nil {
		return "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		if err := os.Rename(nf, filename); err != nil {
			return "", err
		}
		return "", err
	}

	defer func() {
		f.Close()
	}()

	for _, line := range sl {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			return nf, err
		}
	}

	return nf, nil
}
