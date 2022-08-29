package spell

import "html/template"

func renderToken(token string) template.HTML {
	return template.HTML("<input type=\"hidden\" name=\"" + CSRFToken + "\" value=\"" + token + "\" />")
}
