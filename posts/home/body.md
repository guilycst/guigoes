# Home page

## Golang's gopher mascot Midjourney hallucination

![Golang's gopher mascot Midjourney hallucination](./assets/golang_gopher_mascot_17f8cd3e-ac12-4f9d-9ac0-bdaec45b7f84.webp)

> Kinda cute though!

```golang
import (
  "context"
  "fmt"
  "math/rand"
 
  "github.com/OrenRosen/gokit-example/article"
)

type ArticlesRepository interface {
   GetArticle(ctx context.Context, id string) (article.Article, error)
   InsertArticle(ctx context.Context, thing article.Article) error
}

type service struct {
   repo ArticlesRepository
}

func NewService(repo ArticlesRepository) *service {
   return &service{
      repo: repo,
   }
}

func (s *service) GetArticle(ctx context.Context, id string) (article.Article, error) {
   return s.repo.GetArticle(ctx, id)
}

func (s *service) CreateArticle(ctx context.Context, artcle article.Article) (id string, err error) {
   artcle.ID = generateID()
   if err := s.repo.InsertArticle(ctx, artcle); err != nil {
      return "", fmt.Errorf("publishing.CreateArticle: %w", err)
   }
   
   return artcle.ID, nil
}

func generateID() string {
  // code emitted
}
```

# An exhibit of Markdown

This note demonstrates some of what [Markdown][1] is capable of doing.

*Note: Feel free to play with this page. Unlike regular notes, this doesn't automatically save itself.*

## Basic formatting

Paragraphs can be written like so. A paragraph is the basic block of Markdown. A paragraph is what text will turn into when there is no reason it should become anything else.

Paragraphs must be separated by a blank line. Basic formatting of *italics* and **bold** is supported. This *can be **nested** like* so.

## Lists

### Ordered list

1. Item 1
2. A second item
3. Number 3
   1. sub
   2. list
4. Ⅳ
5. und
  * noch
  * mehr
6. weiter

*Note: the fourth item uses the Unicode character for [Roman numeral four][2].*

### Unordered list

* An item
* Another item
* Yet another item
* test
  * test2
  * An item
  * Another item
    * blub
  * narf
* bla

## Checklist

- [ ] a bigger project
  - [ ] first subtask #1234
  - [x] follow up subtask #4321
  - [x] final subtask cc @mention
- [ ] a separate task

## Paragraph modifiers

### Code block

    Code blocks are very useful for developers and other people who look at code or other things that are written in plain text. As you can see, it uses a fixed-width font.

You can also make `inline code` to add code into other things.

```javascript
// This function returns a string padded with leading zeros
function padZeros(num, totalLen) {
   var numStr = num.toString();             // Initialize return value as string
   var numZeros = totalLen - numStr.length; // Calculate no. of zeros
   for (var i = 1; i <= numZeros; i++) {
      numStr = "0" + numStr;
   }
   return numStr;
}
```

## ASCII art

# ○◻◊◇◆■▥▤🗹🗷🗸🗶🖕🕱🌢🏶🕆🕇🕀🖸🖫🗐🖹🗅🖉

<pre>
             ,-. 
    ,     ,-.   ,-. 
   / \   (   )-(   ) 
   \ |  ,.>-(   )-< 
    \|,' (   )-(   ) 
     Y ___`-'   `-' 
     |/__/   `-' 
     | 
     | 
     |    -hrr- 
  ___|_____________ 
</pre>

### Quote

> Here is a quote. What this is should be self explanatory. Quotes are automatically indented when they are used.

## Headings

There are six levels of headings. They correspond with the six levels of HTML headings. You've probably noticed them already in the page. Each level down uses one more hash character.

### Headings *can* also contain **formatting**

### They can even contain `inline code`

Of course, demonstrating what headings look like messes up the structure of the page.

I don't recommend using more than three or four levels of headings here, because, when you're smallest heading isn't too small, and you're largest heading isn't too big, and you want each size up to look noticeably larger and more important, there there are only so many sizes that you can use.

## URLs

URLs can be made in a handful of ways:

* A named link to [MarkItDown][3]. The easiest way to do these is to select what you want to make a link and hit `Ctrl+L`.
* Another named link to [MarkItDown](http://www.markitdown.net/)
* Sometimes you just want a URL like <http://www.markitdown.net/>.

## Horizontal rule

A horizontal rule is a line that goes across the middle of the page.

---

It's sometimes handy for breaking things up.

## Tables

| Name | Description          |
| ------------- | ----------- |
| Help      | ~~Display the~~ help window.|
| Close     | _Closes_ a window     |

Finally, by including colons : within the header row, you can define text to be left-aligned, right-aligned, or center-aligned:

| Left-Aligned  | Center Aligned  | Right Aligned |
| :------------ |:---------------:| -----:|
| col 3 is      | some wordy text | $1600 |
| col 2 is      | centered        |   $12 |
| zebra stripes | are neat        |    $1 |

## Images

![Image of flowers](https://c1.staticflickr.com/1/675/22381378833_c35b2fdd75_n.jpg)

## Finally

There's actually a lot more to Markdown than this :camel:. See the official [introduction][4] and [syntax][5] for more information. However, be aware that this is not using the official implementation, and this might work subtly differently in some of the little things.


  [1]: http://daringfireball.net/projects/markdown/
  [2]: http://www.fileformat.info/info/unicode/char/2163/index.htm
  [3]: http://www.markitdown.net/
  [4]: http://daringfireball.net/projects/markdown/basics
  [5]: http://daringfireball.net/projects/markdown/syntax