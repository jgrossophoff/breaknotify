// breaknotify requires Linux or Darwin with
// either notify-send or osascript installed.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"runtime"
	"time"
)

var (
	endOfNotifications = flag.Duration("end", time.Hour*18, "The time of day to end notifications")
	notifyOnStart      = flag.Bool("s", false, "Show notification on start")
	notifyInterval     = flag.Duration("i", time.Minute*90, "Notification interval")
)

const (
	defaultNotification = `Time to take a break!`
	os                  = runtime.GOOS
)

func main() {
	if os != "linux" && os != "darwin" {
		log.Fatalf("OS %s not supported. Only linux and darwin for now.\n", os)
	}

	if !notificationBinExists() {
		log.Fatalf("make sure that either notify-send or osascript are installed\n")
	}

	flag.Parse()

	if *notifyOnStart {
		notify()
	}

check:
	now := time.Now()
	todaysEnd := time.Date(now.Year(), now.Month(), now.Day(), int(endOfNotifications.Hours()), 0, 0, 0, time.Local)
	if time.Now().After(todaysEnd) {
		// Remind me tomorrow at 9 o'clock
		t := time.Date(now.Year(), now.Month(), now.Day()+1, 8, 0, 0, 0, time.Local)
		time.Sleep(t.Sub(now))
	} else {
		time.Sleep(*notifyInterval)
		notify()
	}
	goto check
}

func notificationBinExists() bool {
	var err error
	if os == "linux" {
		_, err = exec.LookPath("notify-send")
	} else if os == "darwin" {
		_, err = exec.LookPath("osascript")
	}
	return err == nil
}

// notify will call either notify-send or osascript
// to trigger a OS notification.
func notify() {
	msg, err := fetchRandomProverb()

	if os == "linux" {
		if err == nil {
			err := exec.Command("notify-send", msg+"\nSo take a break!").Run()
			if err != nil {
				log.Println(err)
			}
		} else {
			err := exec.Command("notify-send", defaultNotification).Run()
			if err != nil {
				log.Println(err)
			}
		}
	} else if os == "darwin" {
		arg := fmt.Sprintf("display notification \"%s\" with title \"%s\"", msg, defaultNotification)
		err := exec.Command("osascript", "-e", arg).Run()
		if err != nil {
			log.Println(err)
		}
	}
}

// fetchRandomProverb will try to fetch a random proverb
// from http://proverbs-app.antjan.us.
// EDIT service went down so the list is now embedded.
func fetchRandomProverb() (string, error) {
	l := len(proverbs)
	return proverbs[rand.New(rand.NewSource(time.Now().Unix())).Intn(l)-1], nil
}

/*
func fetchRandomProverb() (string, error) {
	var cl http.Client
	cl.Timeout = time.Second * 7
	resp, err := cl.Get("http://proverbs.jens-gross-ophoff.de")
	if err != nil {
		log.Println(err)
		return "", nil
	}

	var prvbs proverbs

	err = json.NewDecoder(resp.Body).Decode(&prvbs)
	if err != nil {
		log.Println(err)
		return "", err
	}

	l := len(prvbs)
	if l == 0 {
		return "", err
	}

	return prvbs[rand.New(rand.NewSource(time.Now().Unix())).Intn(l)-1], nil
}
*/

var proverbs = []string{
	"Dance like nobody is watching, code like everybody is.",
	"A deployed MVP is worth two prototyped.",
	"When you reach bearded-level, there are at least a hundred grey-beards above you.",
	"A/B Test twice, deploy changes once.",
	"Don't commit on master when drunk.",
	"Sleep on a force push.",
	"A git pull a day, keeps the doctor away.",
	"Sometimes you have to cut legacy support to allow the new product to bloom.",
	"More hours worked, more commits made. Mostly reverts and bug-causing features.",
	"Even a greybeard will drop production DB.",
	"Scope creep makes a mountain.",
	"A hundred programmers won't make a two-year project in a week.",
	"Facebook wasn't built in a day.",
	"\"Just ship\" is no substitute for design.",
	"Today's fashion is tomorrow's legacy.",
	"Learning obscure and strange languages, yields better understanding and broader horizons.",
	"The better job you do, the easier others discount the level of difficulty.",
	"Testing is easier than debugging.",
	"Finish a product in a day, and people will expect a new product every day. Teach people about proper development cycles, and your company will flourish.",
	"Customers are the best testers.",
	"Absence is beauty, in error logs.",
	"Eternal sunshine of the stateless mind.",
	"Laziness is your best friend.  Never do twice what you can automate once.",
	"Good test coverage + automated workflows = quiet cell phones and better sleep.",
	"The best code is no code at all.",
	"The best request is the one you don't make.",
	"If a system works perfectly, no one will care what is inside it. Once it breaks, systems design and architecture decides your fate.",
	"Leave architecture for applications that require long-term support.",
	"Architecture and design are preparations for problems and changes, not a key to runtime.",
	"Without a prototype, don't build a final product.",
	"Without boilerplate, there's no speedy development.",
	"Code frustration is a bad advisor for a refactor.",
	"The more technology you learn, the more you realize how little you know.",
	"An early BETA launch will teach you more than a delayed promise.",
	"All applications are pretty when your screen is off.",
	"Do not pick a framework for its demo page, instead pick it for its code.",
	"You cannot set a web standard alone.",
	"A poor programmer blames the language.",
	"The code's writin' but ain't nobody programming.",
	"Mañana often has the most tickets.",
	"Never optimize before measuring",
	"Think about your dance moves when drunk, next time you try to code with some beers on your count.",
	"What happens in Git stays in Git",
	"Simpler code has less bugs.",
	"Lock up your dependency versions and other valuables.",
	"Quantity of attempts often yields quality at the end. Commitment to refactoring legacy code yields better quality yet.",
	"Accept that some days you're the QA and some days you're the one fixing bugs.",
	"Give a programmer the correct code and he can do his work for a day. Teach a programmer to debug and he can do his work for a lifetime - by Chirag Gude",
	"Debugging becomes significantly easier if you first admit that you are the problem.",
	"Figure out your data structures, and the code will follow.",
	"One thing should never do more than one thing.",
	"Success from a final version is a lie, there is only iteration. Through iteration, we gain better products. Through better products, we gain traction. Through traction, we gain success. Through success, misguided tech specs are broken. The development cycle shall free us.",
	"An open source developer does not act for personal fame.",
	"Public code review forces one to better oneself. It forces better practices, smarter solutions, growth as a developer... or being broken.",
	"Testing covers not testing.",
	"The most attractive pull requests are the ones wearing a lot of red.",
	"Coding styleguides without peer code reviews are like running a country on voluntary taxes",
	"Deploying an unmonitored app is like going on a roadtrip without a gas gauge.",
	"Learn a programming language, become a new developer.",
	"Some old code never refactors, and breaks at the slightest change.",
	"A developer will spot a peer from far away",
	"A developer that codes until burnout, lives without a mind.",
	"A marketer is not a QA, a developer does not advertise.",
	"A soft spoken developer will see his warnings of technical debt unheeded, and will suffer the blame.",
	"A well spoken developer can be hired hastily but at the last minute fail an easy test.",
	"One can self-learn the art of code but do not assume other crafts suffer of such low bar of entry.",
	"Find ease in your code: Code difficult to read and understand is code destined to be in troubled legacy.",
	"Collaborating on open source projects can bring about friendship and community just as it can create factions and flame wars.",
	"Refactor or rewrite, there is no patching unmaintainable legacy code",
	"If you stop learning now and take the easy path, you will find yourself stuck in legacy software forever.",
	"A beautiful product which is pleasing to non-paying users is good only for frightening investors when it runs out of funding.",
	"A foreach loop avoided is a CPU cycle earned.",
	"You cannot prevent managers from asking too much of you, but you can prevent them from getting used to it.",
	"Any sufficiently complex app architecture is indistinguishable from spaghetti code.",
	"Writing requirements based code and walking on water are both relatively easy to do when frozen.",
	"It takes twice as much intelligence to debug than to program, therefore you peer review because you can never truly be smart enough to debug your own code.",
	"Hofstadter's Law will tell you to always add more time than you think you need to a project because it will take longer, even when you take into account Hofstadter's Law.",
	"Long lasting code is written only when you pretend that it will be peer reviewed or maintained by a violent psycopath who knows where you live.",
	"Small bug becomes a huge problem.",
	"Commiting is the only command I know, Commiting on you.",
}
