package main

import "strings"

func makeMail(arg3 string) error {
	htmlMail := ds.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
	plainMail := ds.RootPath + "/mail/" + strings.ToLower(arg3) + ".txt.tmpl"

	err := makeFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/mailer/mail.txt.tmpl", plainMail)
	if err != nil {
		return err
	}

	return nil
}
