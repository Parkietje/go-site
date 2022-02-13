# Simple web app written in Go

Using Go `html/template` package for server-side rendering of html templates and [some sprinkles of javascript](https://stimulus.hotwired.dev/) for client-side interactivity.



## Run application
First, create a `.env` file inside `/app` directory:

```
$ cd app/ \
  touch .env \ 
  echo ADMIN="\"""\"" >> .env && \
  echo ADMIN_PASSWORD="\"""\"" >> .env
```

After adding your credentials in `.env`, build the app directory:

`$ go build .`

Finally, exectute the resulting binary:

`$ ./app`

## Templates

The template files, ending in`.gtpl`, inside the `/ui` directory contain plain html, enriched with Go template tags to dynamically generate the UI.


The example below shows the definition of a *base* html template, that is populated with content defined in *title*, *header*, *main*, and *footer* templates. 

![](docs/template_example.png)

Refer to the [Go wiki](https://go.dev/doc/articles/wiki/) for more examples.


## Client-side Interaction

[Stimulus](https://stimulus.hotwired.dev/) is a lightweight javascript library with which you can easily make specific html elements interactive. 

The stimulus controller name for an html component is specified with the *data-controller* tag. The corresponding javascript class will then apply to that element. Inside it you can add elements with *data-actions* tag to trigger methods of the controller class, and *data-targets* tag for elements that we want to control with stimulus.

![](docs/stimulus_example.png)


## Authentication

Users can log in with a username and password. After authentication, session cookies are used to get access to private content.

For now, usernames and passwords are encrypted and stored in files, which are embedded in the binary. The admin password is used as a secret key for encryption/decryption of account data.

## Hot Reload

Install [Air](https://github.com/cosmtrek/air) to enable hot-reloading. After installing, go to app directory, modify `.air.toml`, and run with:

`$ air -c .air.toml `