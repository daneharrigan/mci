package notifier

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/daneharrigan/mci/config"
	"github.com/daneharrigan/mci/db"
)

var (
	cfg     = config.New()
	funcMap = template.FuncMap{"series": seriesHelper}
)

type Notifier interface {
	Send() error
}

func New(user *db.User, comics []*db.Comic) Notifier {
	return &notifier{user: user, comics: comics}
}

type notifier struct {
	user   *db.User
	comics []*db.Comic
}

func (n *notifier) Send() error {
	b, err := n.buf()
	if err != nil {
		return err
	}

	u, err := url.Parse(cfg.MailgunURL)
	if err != nil {
		return err
	}

	u.Path = fmt.Sprintf("/v3/%s/messages", cfg.MailgunDomain)

	q := u.Query()
	q.Set("from", cfg.MailerFrom)
	q.Set("to", n.user.Email)
	q.Set("subject", cfg.MailerSubject)
	q.Set("html", b.String())

	params := bytes.NewBufferString(q.Encode())
	req, err := http.NewRequest("POST", u.String(), params)
	if err != nil {
		return err
	}

	req.SetBasicAuth("api", cfg.MailgunAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("params=%s", q.Encode())
	io.Copy(os.Stdout, res.Body)
	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return res.Body.Close()
}

func (n *notifier) buf() (*bytes.Buffer, error) {
	tmpl, err := template.New("email").Funcs(funcMap).Parse(Template)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	return &b, tmpl.Execute(&b, n.comics)
}

func seriesHelper(comic *db.Comic) string {
	return comic.Series().Name
}
