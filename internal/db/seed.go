package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	repository "github.com/na1tto/go-social/internal/store"
)

var usernames = []string{
	"James", "Mary", "Michael", "Patricia", "John", "Jennifer", "Robert",
	"Linda", "David", "Elizabeth", "William", "Barbara", "Richard", "Susan",
	"Joseph", "Jessica", "Thomas", "Karen", "Christopher", "Sarah", "Charles",
	"Lisa", "Daniel", "Nancy", "Matthew", "Sandra", "Anthony", "Ashley",
	"Mark", "Emily", "Steven", "Kimberly", "Donald", "Betty", "Andrew",
	"Margaret", "Joshua", "Donna", "Paul", "Michelle", "Kenneth", "Carol",
	"Kevin", "Amanda", "Brian", "Melissa", "Timothy", "Deborah", "Ronald",
	"Stephanie",
}

var titles = []string{
	"5 Habits of Highly Effective Coders",
	"The Minimalist's Guide to Digital Declutter",
	"Is AI Taking Over Your Job? A Reality Check",
	"One-Pot Meals for the Busy Weeknight",
	"Mastering the Command Line in 15 Minutes",
	"Boost Your Focus with the Pomodoro Technique",
	"The Secret to Writing Better, Faster Emails",
	"Affordable Home Decor Upgrades That Work",
	"Understanding Go Routines: Simply Explained",
	"How to Start Investing with Only $100",
	"My Favorite Free Tools for Web Development",
	"Why You Need a Digital Detox This Weekend",
	"The Best Budget Laptops of the Year",
	"Simple Steps to a Healthier Morning Routine",
	"Ditch the Distractions: Ultimate Productivity Hacks",
	"Testing Your Go Code: A Quick Start Guide",
	"Travel Hacking 101: Fly Cheaper, Stay Longer",
	"The Power of Saying 'No' (and How to Do It)",
	"Quick Fixes for Slow Wi-Fi Speeds",
	"Beyond the To-Do List: The Next Level of Planning",
}

var contents = []string{
	"Great coding is often less about raw talent and more about consistent, disciplined habits. We dive into five key practices—from effective debugging to continuous learning—that separate the good developers from the truly exceptional ones.",
	"Our digital lives are often as cluttered as our physical ones. This guide walks you through the essential steps for cleaning up your desktop, streamlining your inbox, and freeing up storage space for a calmer, more focused workflow.",
	"The rise of generative AI has sparked both excitement and fear in the job market. Instead of panicking, let's look at the data: which roles are truly at risk, and more importantly, how can you leverage AI to make your own work indispensable?",
	"Tired of washing multiple pans after a long day? One-pot cooking is the ultimate solution for quick, healthy, and low-mess dinners. Here are our top five recipes that minimize cleanup without sacrificing flavor.",
	"The command line can seem intimidating, but its power and efficiency are undeniable. We break down the absolute must-know commands and shortcuts that will elevate your productivity in just a quarter of an hour.",
	"Distraction is the enemy of productivity. The Pomodoro Technique, which uses timed intervals of intense work and short breaks, is a simple yet powerful method to keep your focus sharp and your burnout low. Learn how to implement it today.",
	"Email is a necessary evil, but it doesn't have to consume your day. By adopting a few simple principles—like front-loading key information and using clear calls to action—you can write professional emails in half the time.",
	"You don't need a massive budget to transform your living space. We share brilliant, low-cost ideas—from swapping out fixtures to mastering lighting—that provide maximum visual impact for minimal financial outlay.",
	"Concurrency is one of Go's greatest strengths, and goroutines are at its heart. This article strips away the complex jargon to explain exactly what a goroutine is, how it works, and why it makes Go programs so incredibly efficient.",
	"Many people delay investing because they think they need a large sum of money. We prove that wrong by detailing exactly how and where you can begin building your portfolio right now, even if you only have a Benjamin to start with.",
	"Building modern websites requires a stack of tools, but not all of them need to cost a fortune. From code editors to testing platforms, here is a curated list of high-quality, free resources that every developer should be using.",
	"Constant connection leads to mental fatigue and reduced creativity. A weekend digital detox is a powerful reset button for your brain. We provide a step-by-step plan to unplug successfully and reconnect with the physical world.",
	"Finding a powerful, reliable laptop doesn't have to drain your bank account. Our annual roundup identifies the top performers in the budget category, analyzing speed, battery life, and build quality to help you choose wisely.",
	"How you start your day often dictates how the rest of it unfolds. Forget the hour-long yoga sessions—we focus on small, sustainable changes you can implement immediately to make your mornings more energized and productive.",
	"From push notifications to open office plans, distractions are everywhere. This guide offers advanced strategies and unconventional hacks—beyond just turning off your phone—to create a deep-focus environment and maximize output.",
	"Writing effective tests is crucial for maintaining reliable Go applications. We provide a rapid walkthrough of Go's built-in `testing` package, showing you how to write your first unit tests and simple benchmarks in minutes.",
	"Travel hacking isn't just for experts. Learn the fundamental strategies for accumulating points, finding error fares, and leveraging credit card rewards that will allow you to see the world without spending a fortune.",
	"Overcommitment is the silent killer of personal time and productivity. Mastering the art of the polite but firm 'no' is a game-changer. This article provides practical scripts and mindset shifts for protecting your most valuable resource: your time.",
	"A slow internet connection can bring your productivity to a halt. Before calling your provider, try these five simple troubleshooting steps—from optimizing router placement to clearing network caches—that often resolve lag instantly.",
	"The classic to-do list is often overwhelming and ineffective. We explore alternative and modern planning methodologies, such as time blocking and the Ivy Lee method, that help you prioritize and actually finish your most important tasks.",
}

var tags = []string{
	"Productivity",
	"GoLang",
	"WebDevelopment",
	"AI",
	"Career",
	"CodingTips",
	"PersonalFinance",
	"Investing",
	"HomeImprovement",
	"Minimalism",
	"TravelHacks",
	"Cooking",
	"TimeManagement",
	"Networking",
	"TechReview",
	"Health",
	"Routine",
	"SoftwareTesting",
	"Communication",
	"LifeHacks",
}

var comments = []string{
	"Great post! Thanks for sharing these informations!",
	"I completlely agree with your thoughts.",
	"Thanks for the tips, very helpful!",
	"Interesting perspective, I hadn't considered that.",
	"Really useful information right here, thank you.",
	"Very interesting, I'm looking foward for your next posts!",
	"That's very inspirational, thank you!",
	"Very insightful post",
	"What a way of thinking about this subject, really insightful.",
	"I was skeptical at the begninning but after reading this post it changed my opnion on this subject.",
	"Well written, I enjoyed reading this.",
}

// this is the script for seeding the database with test data

func Seed(store repository.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
		}
	}

	log.Println("The seeding was completed!")
}

func generateUsers(num int) []*repository.User {
	users := make([]*repository.User, num)

	for i := 0; i < num; i++ {
		users[i] = &repository.User{
			UserName: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Emai:     usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*repository.User) []*repository.Post {
	posts := make([]*repository.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &repository.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(titles))],
			Tags: []string{
				tags[rand.Intn(len(titles))],
				tags[rand.Intn(len(titles))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*repository.User, posts []*repository.Post) []*repository.Comment {
	cms := make([]*repository.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &repository.Comment{
			PostId:  posts[rand.Intn(len(posts))].ID,
			UserId:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
