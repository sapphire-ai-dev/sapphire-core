package main

import (
	"context"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	chromedp.Run(ctx, chromedp.Navigate("https://www.google.com"))
}
