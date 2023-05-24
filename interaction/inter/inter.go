package inter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/h3ll0kitt1/tg_bot_for_indecisive/https/telegram"
	"github.com/h3ll0kitt1/tg_bot_for_indecisive/storage"
)

const (
	cmdHelp     = "/help"
	cmdSave     = "/save"
	cmdDone     = "/done"
	cmdSurprise = "/surprise"
	cmdPrint    = "/list"
)

func Process(storage storage.Stored, in telegram.UpdatesResponse) (*telegram.UpdatesResponse, error) {
	switch {
	case checkCmd(&in, cmdSave):
		return save(storage, &in)
	case checkCmd(&in, cmdDone):
		return done(storage, &in)
	case checkCmd(&in, cmdSurprise):
		return surprise(storage, &in)
	case checkCmd(&in, cmdPrint):
		return print(storage, &in)
	default:
		return help(&in), nil
	}
	return &in, nil
}

func checkCmd(in *telegram.UpdatesResponse, cmd string) bool {
	if !strings.Contains(in.Result[0].Message.Text, cmd) {
		return false
	}
	return true
}

func book(u *telegram.UpdatesResponse) string {
	return strings.Join(strings.Split(u.Result[0].Message.Text, " ")[1:], " ")
}

func chat(u *telegram.UpdatesResponse) int {
	return u.Result[0].Message.Chat.Id
}

func updateText(u *telegram.UpdatesResponse, newText string) {
	u.Result[0].Message.Text = newText
}

func validate(in string) bool {
	isStringAlphabetic := regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
	if len(in) == 0 || !isStringAlphabetic(in) {
		return false
	}
	return true
}

func save(storage storage.Stored, in *telegram.UpdatesResponse) (*telegram.UpdatesResponse, error) {
	if !validate(book(in)) {
		message := "Book name should only contain letters or digits\n"
		updateText(in, message)
		return in, nil
	}

	ok, err := storage.Save(chat(in), book(in))
	if err != nil {
		return nil, fmt.Errorf("Save to storage: %w", err)
	}

	message := book(in) + " is saved to your list\n"
	if !ok {
		message = book(in) + " is already in your list\n"
	}
	updateText(in, message)
	return in, nil
}

func done(storage storage.Stored, in *telegram.UpdatesResponse) (*telegram.UpdatesResponse, error) {
	if !validate(book(in)) {
		message := "Book name should only contain letters or digits\n"
		updateText(in, message)
		return in, nil
	}

	ok, err := storage.Delete(chat(in), book(in))
	if err != nil {
		return nil, fmt.Errorf("Delete from storage: %w", err)
	}

	message := book(in) + " is deleted from your list\n"
	if !ok {
		message = book(in) + " has not been in your list\n"
	}
	updateText(in, message)
	return in, nil
}

func print(storage storage.Stored, in *telegram.UpdatesResponse) (*telegram.UpdatesResponse, error) {
	ok, err := storage.LenNotZero(chat(in))
	if err != nil {
		return nil, err
	}

	message := "Your list is empty\n"
	if ok {
		books, err := storage.Print(chat(in))
		if err != nil {
			return nil, fmt.Errorf("Print book list: %w", err)
		}
		message = strings.Join(books, "\n")
	}
	updateText(in, message)
	return in, nil
}

func surprise(storage storage.Stored, in *telegram.UpdatesResponse) (*telegram.UpdatesResponse, error) {
	ok, err := storage.LenNotZero(chat(in))
	if err != nil {
		return nil, fmt.Errorf("Check list length: %w", err)
	}

	message := ""
	if !ok {
		message = "Your list is empty, please, add more books\n"
	} else {
		rnd, err := storage.Rand(chat(in))
		if err != nil {
			return nil, fmt.Errorf("Get ramdon book from list: %w", err)
		}

		message = "You should read " + rnd + "\n"
	}
	updateText(in, message)
	return in, nil
}

func help(in *telegram.UpdatesResponse) *telegram.UpdatesResponse {

	message := "This Bot helps indecisisve people to choose books from their large <should read list>.\n" +
		"\nPlease enter:\n" +
		"/save <book_name> - to save book to list\n" +
		"/done <book_name> - to delete book from lists\n" +
		"/surpise - to ask Bot to make decision for you\n" +
		"/list - to print your current list\n"

	updateText(in, message)
	return in
}
