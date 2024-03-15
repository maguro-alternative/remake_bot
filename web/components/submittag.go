package components

func CreateSubmitTag(permission string) string {
	var submitTag string
	if permission == "write" || permission == "all" {
		submitTag = `<button type="submit" class="btn btn-primary">送信</button>`
	}
	return submitTag
}