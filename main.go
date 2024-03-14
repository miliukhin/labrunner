package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"os/exec"
	"github.com/go-rod/rod"
)

func StartNextAttempt(page *rod.Page) error {
	button, err := page.MustWaitLoad().Timeout(time.Second).ElementR("button", "Зробити наступну спробу")
	if err != nil {
		button, err = page.MustWaitLoad().Timeout(time.Second).ElementR("button", "Спроба тесту")
		if err != nil {
			button, err = page.MustWaitLoad().Timeout(time.Second).ElementR("button", "Продовжуйте свою спробу")
		}
	}
	button.MustClick()
	return nil
}

func FinishTest(page *rod.Page) (error, bool) {
	page.MustWaitLoad()
	button, err := page.Timeout(time.Second).Element("input[type='submit'][value='Наступна сторінка']")
	if err != nil {
		button := page.MustElement("input[type='submit'][value='Завершити спробу...']")
		button = page.MustWaitLoad().MustElementR("button", "Відправити все та завершити")
		modal := page.MustElement(".modal-footer")
		button = modal.MustElementR("button", "Відправити все та завершити")
		fmt.Printf("%s\n", string(button.MustText()))
		return nil, true
	}
	button.MustClick()
	return nil, false
}

// Run the question text and put results in a corresponding form
func RunTasks(page *rod.Page) error {
	page.MustWaitLoad()
	tests := page.MustElements(".formulation.clearfix")
	for _, element := range tests {
		Code := element.MustElement(".qtext").MustText()
		// Code := element.MustElementR("p", "print()").MustText()
		fmt.Println(Code)
		result, err := exec.Command("python", "-c", Code).Output()
		fmt.Printf("%s\n", string(result))
		if err != nil {
			// log.Print(err)
			fmt.Printf("%s\n", err)
		} else {
			form, formerr := element.Element(".form-control")
			radioBoxes, raderr := element.Elements("input[type='radio']")
			if formerr == nil {
	        		form.MustInput(string(result))
			} else if raderr == nil {
				for len(result) > 0 && (result[len(result)-1] == '\n' || result[len(result)-1] == '\r') {
        				result = result[:len(result)-1]
    				}
				// fmt.Println("solving radio:")
				answers := element.MustElements(".flex-fill.ml-1")
				for i, a := range answers {
					// fmt.Println(string(result))
					// fmt.Println(a.MustText())
					if a.MustText() == string(result) {
						radioBoxes[i].MustClick()
						break
					}
				}
			}
		}
	}
	return nil
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("name <link> <email> <password>")
		return
	}
	link	:= os.Args[1]
	email	:= os.Args[2]
	pass	:= os.Args[3]

	page := rod.New().NoDefaultDevice().MustConnect().MustPage(link)
	page.MustElement("#region-main > div > div > div > div > div.login-identityproviders > a").MustClick()
	page.MustElement("#identifierId").MustInput(email)
	page.MustElement("#identifierNext > div > button").MustClick()
	page.MustElement("#password").MustInput(pass)
	page.MustElement("#passwordNext > div > button").MustClick()

	err := StartNextAttempt(page)
	if err != nil {
		log.Fatalln("Error while starting new attempt:", err)
	}
	isLastPage := false
	for !isLastPage {
		err = RunTasks(page)
		if err != nil {
			log.Fatalln("Error while making test: ", err)
		}
		_, isLastPage = FinishTest(page)
	}

	time.Sleep(time.Hour)
}
