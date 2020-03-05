package main

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "assets/main.css",
		FileModTime: time.Unix(1583314864, 0),

		Content: string("html,\nbody {\n    height: 100%;\n}\n\nbody {\n    display: -ms-flexbox;\n    display: flex;\n    -ms-flex-align: center;\n    align-items: center;\n    padding-top: 40px;\n    padding-bottom: 40px;\n    background-color: #f5f5f5;\n}\n\n.form-signin {\n    width: 100%;\n    max-width: 420px;\n    padding: 15px;\n    margin: auto;\n}\n\n.form-label-group {\n    position: relative;\n    margin-bottom: 1rem;\n}\n\n.form-label-group > input,\n.form-label-group > label {\n    height: 3.125rem;\n    padding: .75rem;\n}\n\n.form-label-group > label {\n    position: absolute;\n    top: 0;\n    left: 0;\n    display: block;\n    width: 100%;\n    margin-bottom: 0; /* Override default `<label>` margin */\n    line-height: 1.5;\n    color: #495057;\n    pointer-events: none;\n    cursor: text; /* Match the input under the label */\n    border: 1px solid transparent;\n    border-radius: .25rem;\n    transition: all .1s ease-in-out;\n}\n\n.form-label-group input::-webkit-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input:-ms-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::-ms-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::-moz-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::placeholder {\n    color: transparent;\n}\n\n.form-label-group input:not(:placeholder-shown) {\n    padding-top: 1.25rem;\n    padding-bottom: .25rem;\n}\n\n.form-label-group input:not(:placeholder-shown) ~ label {\n    padding-top: .25rem;\n    padding-bottom: .25rem;\n    font-size: 12px;\n    color: #777;\n}\n\n/* Fallback for Edge\n-------------------------------------------------- */\n@supports (-ms-ime-align: auto) {\n    .form-label-group > label {\n        display: none;\n    }\n    .form-label-group input::-ms-input-placeholder {\n        color: #777;\n    }\n}\n\n/* Fallback for IE\n-------------------------------------------------- */\n@media all and (-ms-high-contrast: none), (-ms-high-contrast: active) {\n    .form-label-group > label {\n        display: none;\n    }\n    .form-label-group input:-ms-input-placeholder {\n        color: #777;\n    }\n}"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "assets/main_old.css",
		FileModTime: time.Unix(1583302866, 0),

		Content: string(".login-form {\n    width: 340px;\n    margin: 50px auto;\n}\n\n.login-form form {\n    margin-bottom: 15px;\n    background: #f7f7f7;\n    box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);\n    padding: 30px;\n}\n\n.login-form h2 {\n    margin: 0 0 15px;\n}\n\n.form-control,\n.btn {\n    min-height: 38px;\n    border-radius: 2px;\n}\n\n.btn {\n    font-size: 15px;\n    font-weight: bold;\n}\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "flashes.html",
		FileModTime: time.Unix(1583339840, 0),

		Content: string("{{ if .error }}\n    <div class=\"alert alert-danger\" role=\"alert\">\n        {{.error}}\n    </div>\n{{ end }}\n\n{{ if .message }}\n    <div class=\"alert alert-primary\" role=\"alert\">\n        {{.message}}\n    </div>\n{{end}}"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "forget.html",
		FileModTime: time.Unix(1583382088, 0),

		Content: string("{{ define \"head\" }}\n<title>Forgot your password?</title>\n{{end}}\n\n{{define \"content\"}}\n    <form class=\"form-signin\" action=\"\" method=\"post\">\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Forgot your password?</h1>\n        </div>\n\n        {{ include \"flashes\" }}\n        <div class=\"form-label-group\">\n            <input type=\"email\" id=\"inputEmail\" name=\"email\" autocomplete=\"email\" class=\"form-control\"\n                   placeholder=\"Email\" required\n                   autofocus>\n            <label for=\"inputEmail\">Email</label>\n        </div>\n\n        <button class=\"btn btn-lg btn-primary btn-block\" type=\"submit\">Recover my account</button>\n        <p class=\"text-center p-2\"><a href=\"/login\">I have an account</a></p>\n\n\n    </form>\n{{end}}"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "index.html",
		FileModTime: time.Unix(1583338799, 0),

		Content: string("{{ define \"head\"}}\n    <title>Welcome</title>\n{{end}}\n{{ define \"content\"}}\n\n    <form class=\"form-signin\" action=\"\" method=\"post\">\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Welcome</h1>\n        </div>\n\n        {{ include \"flashes\" }}\n        <div class=\"form-label-group\">\n            <input type=\"text\" id=\"inputUsername\" name=\"username\" autocomplete=\"username\" class=\"form-control\"\n                   placeholder=\"Username\" required\n                   autofocus>\n            <label for=\"inputUsername\">Username</label>\n        </div>\n\n        <div class=\"form-label-group\">\n            <input type=\"password\" autocomplete=\"current-password\" name=\"password\" id=\"inputPassword\"\n                   class=\"form-control\" placeholder=\"Password\"\n                   required>\n            <label for=\"inputPassword\">Password</label>\n        </div>\n\n        <div class=\"checkbox\">\n            <label>\n                <input type=\"checkbox\" name=\"remember\" value=\"remember-me\"> Remember me\n            </label>\n        </div>\n\n        <div class=\"mb-3 mt-1\">\n            <a href=\"/forget\" class=\"\">Forgot Password?</a>\n        </div>\n\n        <button class=\"btn btn-lg btn-primary btn-block\" type=\"submit\">Sign in</button>\n        <p class=\"text-center p-2\"><a href=\"/register\">Create an Account</a></p>\n\n\n    </form>\n{{end}}"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "layouts/master.html",
		FileModTime: time.Unix(1583337284, 0),

		Content: string("<!DOCTYPE html>\n<html lang=\"en\">\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n\n    <!-- Bootstrap core CSS -->\n    <link rel=\"stylesheet\" href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css\"\n          integrity=\"sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh\" crossorigin=\"anonymous\">\n\n\n    <link rel=\"stylesheet\" type=\"text/css\" href=\"../assets/main.css\">\n\n\n    {{template \"head\" .}}\n</head>\n\n<body>\n{{template \"content\" .}}\n\n<script src=\"https://cdn.jsdelivr.net/npm/fingerprintjs2@2.1.0/dist/fingerprint2.min.js\"\n        integrity=\"sha384-UxnJUeaUHFTScXFSpWAeq3BsPELZB9qc8o37yeXwfLcAfOhSoirtOIyu4k6GU9lk\"\n        crossorigin=\"anonymous\"></script>\n<script>\n    if (window.requestIdleCallback) {\n        requestIdleCallback(function() {\n            Fingerprint2.get(function(components) {\n                var values = components.map(function (component) { return component.value })\n                var murmur = Fingerprint2.x64hash128(values.join(''), 31)\n                document.getElementById(\"browserfingerprint\").value = murmur;\n                try {\n                    document.getElementById(\"browserfingerprint\").value = murmur;\n                }catch(e) {}\n            })\n        })\n    } else {\n        setTimeout(function() {\n            Fingerprint2.get(function(components) {\n                var values = components.map(function (component) { return component.value })\n                var murmur = Fingerprint2.x64hash128(values.join(''), 31)\n                try {\n                    document.getElementById(\"browserfingerprint\").value = murmur;\n                }catch(e) {}\n            });\n        }, 500);\n    }\n</script>\n</body>\n\n</html>"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "mfa.html",
		FileModTime: time.Unix(1583397649, 0),

		Content: string("{{ define \"head\" }}\n    <title>Security check</title>\n{{end}}\n\n{{define \"content\"}}\n    <form class=\"form-signin\" action=\"\" method=\"post\">\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Security check</h1>\n            <p>Welcome, {{.first_name}}. Please check your email for your code.</p>\n        </div>\n\n        {{ include \"flashes\" }}\n        <div class=\"form-label-group\">\n            <input type=\"text\" inputmode=\"number\" id=\"inputCode\" name=\"code\" class=\"form-control\"\n                   placeholder=\"Code\" required\n                   autofocus>\n            <label for=\"inputCode\">Code</label>\n        </div>\n\n        <button class=\"btn btn-lg btn-primary btn-block\" type=\"submit\">Continue</button>\n        <p class=\"text-center p-2\"><a href=\"/logout\">Logout</a></p>\n\n\n    </form>\n{{end}}"),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "newpassword.html",
		FileModTime: time.Unix(1583394996, 0),

		Content: string("{{ define \"head\"}}\n    <title>Welcome</title>\n{{end}}\n{{ define \"content\"}}\n\n    <form class=\"form-signin\" action=\"\" method=\"post\">\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Updating Password</h1>\n        </div>\n\n        {{ include \"flashes\" }}\n        <div class=\"form-label-group\">\n            <input type=\"password\" autocomplete=\"new-password\" name=\"password\" id=\"inputPassword1\"\n                   class=\"form-control\" placeholder=\"Password\"\n                   required>\n            <label for=\"inputPassword1\">Password</label>\n        </div>\n\n        <div class=\"form-label-group\">\n            <input type=\"password\" autocomplete=\"new-password\" name=\"confirm_password\" id=\"inputPassword2\"\n                   class=\"form-control\" placeholder=\"Confirm Password\"\n                   required>\n            <label for=\"inputPassword2\">Confirm Password</label>\n        </div>\n\n        <button class=\"btn btn-lg btn-primary btn-block\" type=\"submit\">Update password</button>\n\n\n    </form>\n{{end}}"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "register.html",
		FileModTime: time.Unix(1583339979, 0),

		Content: string("{{ define \"head\" }}\n    <title>Registration</title>\n    <script src=\"https://www.google.com/recaptcha/api.js?render={{.recaptchaSite}}\"></script>\n{{end}}\n\n{{ define \"content\" }}\n\n    <form class=\"form-signin\" action=\"/register\" method=\"post\">\n\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Register an Account</h1>\n        </div>\n\n\n        {{ include \"flashes\" }}\n\n\n        <div class=\"form-row\">\n            <div class=\"col\">\n                <div class=\"form-label-group\">\n                    <input type=\"text\" id=\"inputFirstName\" name=\"first_name\" class=\"form-control\" placeholder=\"First Name\" required\n                           autofocus>\n                    <label for=\"inputFirstName\">First Name</label>\n                </div>\n            </div>\n            <div class=\"col\">\n                <div class=\"form-label-group\">\n                    <input type=\"text\" id=\"inputLastName\" name=\"last_name\" class=\"form-control\" placeholder=\"Last Name\" required\n                           >\n                    <label for=\"inputLastName\">Last Name</label>\n                </div>\n            </div>\n        </div>\n\n\n        <div class=\"form-label-group\">\n            <input type=\"text\" id=\"inputUsername\" autocomplete=\"username\" name=\"username\" class=\"form-control\" placeholder=\"Username\" required\n                   >\n            <label for=\"inputUsername\">Username</label>\n        </div>\n\n        <div class=\"form-label-group\">\n            <input type=\"email\" id=\"inputEmail\" autocomplete=\"email\" name=\"email\" class=\"form-control\" placeholder=\"Email\" required\n                   >\n            <label for=\"inputEmail\">Email</label>\n        </div>\n        <div class=\"form-label-group\">\n            <input type=\"password\" autocomplete=\"new-password\" id=\"inputPassword1\" name=\"password\" class=\"form-control\" placeholder=\"Password\" required\n                   >\n            <label for=\"inputPassword1\">Password</label>\n        </div>\n        <div class=\"form-label-group\">\n            <input type=\"password\" id=\"inputPassword2\" autocomplete=\"new-password\" name=\"confirm_password\" class=\"form-control\" placeholder=\"Confirm Password\" required\n                   >\n            <label for=\"inputPassword2\">Confirm Password</label>\n        </div>\n\n        <button class=\"btn btn-lg btn-primary btn-block\" type=\"submit\">Register</button>\n        <p class=\"text-center p-2\"><a href=\"/login\">I have an account</a></p>\n\n\n        <div>\n            <script>\n                grecaptcha.ready(function () {\n                    grecaptcha.execute('{{.recaptchaSite}}', {action: 'register'}).then(function (token) {\n\n                    });\n                });\n            </script>\n        </div>\n\n    </form>\n{{end}}"),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "welcome.html",
		FileModTime: time.Unix(1583395092, 0),

		Content: string("{{ define \"head\"}}\n    <title>Welcome</title>\n{{end}}\n{{ define \"content\"}}\n\n    <div class=\"form-signin\">\n\n        <input type=\"hidden\" name=\"__csrf\" value=\"{{.csrf}}\"/>\n        <input type=\"hidden\" id=\"browserfingerprint\" name=\"browser_fingerprint\" value=\"{{.csrf}}\"/>\n        <div class=\"text-center mb-4\">\n            <h1 class=\"h3 mb-3 font-weight-normal\">Welcome, {{.User.FirstName}}</h1>\n        </div>\n\n        <div>\n\n        </div>\n\n        <a href=\"/logout\" class=\"btn btn-lg btn-primary btn-block\">Sign out</a>\n\n    </div>\n{{end}}"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1583397649, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file5, // "flashes.html"
			file6, // "forget.html"
			file7, // "index.html"
			filea, // "mfa.html"
			fileb, // "newpassword.html"
			filec, // "register.html"
			filed, // "welcome.html"

		},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "assets",
		DirModTime: time.Unix(1583314864, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "assets/main.css"
			file4, // "assets/main_old.css"

		},
	}
	dir8 := &embedded.EmbeddedDir{
		Filename:   "layouts",
		DirModTime: time.Unix(1583337284, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file9, // "layouts/master.html"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "assets"
		dir8, // "layouts"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}
	dir8.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./project`, &embedded.EmbeddedBox{
		Name: `./project`,
		Time: time.Unix(1583397649, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":        dir1,
			"assets":  dir2,
			"layouts": dir8,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"assets/main.css":     file3,
			"assets/main_old.css": file4,
			"flashes.html":        file5,
			"forget.html":         file6,
			"index.html":          file7,
			"layouts/master.html": file9,
			"mfa.html":            filea,
			"newpassword.html":    fileb,
			"register.html":       filec,
			"welcome.html":        filed,
		},
	})
}

func init() {

	// define files
	filef := &embedded.EmbeddedFile{
		Filename:    "main.css",
		FileModTime: time.Unix(1583314864, 0),

		Content: string("html,\nbody {\n    height: 100%;\n}\n\nbody {\n    display: -ms-flexbox;\n    display: flex;\n    -ms-flex-align: center;\n    align-items: center;\n    padding-top: 40px;\n    padding-bottom: 40px;\n    background-color: #f5f5f5;\n}\n\n.form-signin {\n    width: 100%;\n    max-width: 420px;\n    padding: 15px;\n    margin: auto;\n}\n\n.form-label-group {\n    position: relative;\n    margin-bottom: 1rem;\n}\n\n.form-label-group > input,\n.form-label-group > label {\n    height: 3.125rem;\n    padding: .75rem;\n}\n\n.form-label-group > label {\n    position: absolute;\n    top: 0;\n    left: 0;\n    display: block;\n    width: 100%;\n    margin-bottom: 0; /* Override default `<label>` margin */\n    line-height: 1.5;\n    color: #495057;\n    pointer-events: none;\n    cursor: text; /* Match the input under the label */\n    border: 1px solid transparent;\n    border-radius: .25rem;\n    transition: all .1s ease-in-out;\n}\n\n.form-label-group input::-webkit-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input:-ms-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::-ms-input-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::-moz-placeholder {\n    color: transparent;\n}\n\n.form-label-group input::placeholder {\n    color: transparent;\n}\n\n.form-label-group input:not(:placeholder-shown) {\n    padding-top: 1.25rem;\n    padding-bottom: .25rem;\n}\n\n.form-label-group input:not(:placeholder-shown) ~ label {\n    padding-top: .25rem;\n    padding-bottom: .25rem;\n    font-size: 12px;\n    color: #777;\n}\n\n/* Fallback for Edge\n-------------------------------------------------- */\n@supports (-ms-ime-align: auto) {\n    .form-label-group > label {\n        display: none;\n    }\n    .form-label-group input::-ms-input-placeholder {\n        color: #777;\n    }\n}\n\n/* Fallback for IE\n-------------------------------------------------- */\n@media all and (-ms-high-contrast: none), (-ms-high-contrast: active) {\n    .form-label-group > label {\n        display: none;\n    }\n    .form-label-group input:-ms-input-placeholder {\n        color: #777;\n    }\n}"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "main_old.css",
		FileModTime: time.Unix(1583302866, 0),

		Content: string(".login-form {\n    width: 340px;\n    margin: 50px auto;\n}\n\n.login-form form {\n    margin-bottom: 15px;\n    background: #f7f7f7;\n    box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);\n    padding: 30px;\n}\n\n.login-form h2 {\n    margin: 0 0 15px;\n}\n\n.form-control,\n.btn {\n    min-height: 38px;\n    border-radius: 2px;\n}\n\n.btn {\n    font-size: 15px;\n    font-weight: bold;\n}\n"),
	}

	// define dirs
	dire := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1583314864, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filef, // "main.css"
			fileg, // "main_old.css"

		},
	}

	// link ChildDirs
	dire.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./project/assets`, &embedded.EmbeddedBox{
		Name: `./project/assets`,
		Time: time.Unix(1583314864, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dire,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"main.css":     filef,
			"main_old.css": fileg,
		},
	})
}
