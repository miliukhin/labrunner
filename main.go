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
	        	element.MustElement(".form-control").MustInput(string(result))
		}
	}
	return nil
}

// тре субміт тільки тим елементам, де нема Правильно/неправильно + по одному
// так не паше, тоді списку не треба, просто перший без обирається
// але будуть нюанси з початковими тестами? ні, якщо неправильно теж вважати то норм
// кароче не хоче вміти
func SubmitAnswersInRust (page *rod.Page) error {
	elements := page.MustElements(".formulation.clearfix")
	for _, element := range elements {
		// badge := element.MustElementR("div", "равильно")
		badge := element.MustElementX("//div[contains@class='badge']")
		// fmt.Println(badge.MustText())
		if badge == nil {
			fmt.Println("НЄ бейдж")
		} else {
			fmt.Println("Є")
			// button := element.MustElement("input[type='submit'][value='Перевірити']")
			// if button != nil {
			// 	fmt.Println(button.MustText())
			// 	button.MustClick()
			// } else {
			// }
		}
	}
// <div class="correctness badge correct">Правильно</div>
// #yui_3_18_1_1_1710351328609_408 > div:nth-child(1)
// /html/body/div[1]/div[2]/div/div[1]/section/div[2]/form/div/div[6]/div[2]/div[2]/div/div[1]
	return nil
}

// це стара теж погана
func SubmitAnswers (page *rod.Page) error {
	buttons := page.MustElements("input[type='submit'][value='Перевірити']")
	for _, button := range buttons {
		if button != nil {
			button.MustClick()
		} else {
			fmt.Println("Button not found")
		}
	}
	return nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("name <email> <password>")
		return
	}
	email := os.Args[1]
	pass := os.Args[2]

	page := rod.New().NoDefaultDevice().MustConnect().MustPage("https://vns.lpnu.ua/mod/quiz/view.php?id=835577")
	page.MustElement("#region-main > div > div > div > div > div.login-identityproviders > a").MustClick()
	page.MustElement("#identifierId").MustInput(email)
	page.MustElement("#identifierNext > div > button").MustClick()
	page.MustElement("#password").MustInput(pass)
	page.MustElement("#passwordNext > div > button").MustClick()

	// StartNextAttempt(page)
	// SubmitAnswersInRust(page)
	// SubmitAnswers(page)

	err := StartNextAttempt(page)
	if err != nil {
		log.Fatalln("Error while startin new attempt:", err)
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
