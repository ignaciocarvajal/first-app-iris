package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	v1 := app.Party("/")
	{
		v1.Get("/", initPage)
		//v1.Post("/send-email", sendEmail)
	}
	v1.HandleDir("/static", iris.Dir("./assets"), iris.DirOptions{
		// Defaults to "/index.html", if request path is ending with **/*/$IndexName
		// then it redirects to **/*(/) which another handler is handling it,
		// that another handler, called index handler, is auto-registered by the framework
		// if end developer do not managed to handle it by hand.
		IndexName: "/index.html",
		// When files home served under compression.
		Compress: false,
		// List the files inside the current requested directory if `IndexName` not found.
		ShowList: false,
		// When ShowList is true you can configure if you want to show or hide  files.
		ShowHidden: false,
		Cache: iris.DirCacheOptions{
			// enable in-memory cache and pre-compress the files.
			Enable: true,
			// ignore image types (and pdf).
			CompressIgnore: iris.MatchImagesAssets,
			// do not compress files smaller than size.
			CompressMinSize: 300,
			// available encodings that will be negotiated with client's needs.
			Encodings: []string{"gzip", "br" /* you can also add: deflate, snappy */},
		},
		DirList: iris.DirListRich(),
		// If `ShowList` is true then this function will be used instead of the default
		// one to show the list of files of a current requested directory(dir).
		// DirList: func(ctx iris.Context, dirName string, dir http.File) error { ... }
		//
		// Optional validator that loops through each requested resource.
		// AssetValidator:  func(ctx iris.Context, name string) bool { ... }
	})
	//v1.Get("/", func(ctx iris.Context) {
	//	ctx.ServeFile("./assets/index.html")
	//})

	api := app.Party("/api")
	{
		api.Post("/send-email", sendEmail)
	}

	app.Logger().SetLevel("debug")
	err := app.Listen(":9090")
	if err != nil {
		panic(err)
	}
}

func initPage(ctx iris.Context) {
	err := ctx.View("index.html")
	if err != nil {
		ctx.View("404.html")
	}

}

type ContactDetail struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func sendEmail(ctx iris.Context) {
	contact := &ContactDetail{
		Name:    ctx.PostValue("name"),
		Email:   ctx.PostValue("email"),
		Message: ctx.PostValue("message"),
		Subject: ctx.PostValue("subject"),
	}

	fmt.Println("data:", contact)

	from := mail.NewEmail(contact.Name, "contact@ignaciocarvajald.us")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Contacto", "contact@ignaciocarvajald.us")
	plainTextContent := contact.Message
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		_, err = ctx.Text("OK")
		if err != nil {
			panic(err)
		}
	}

}
