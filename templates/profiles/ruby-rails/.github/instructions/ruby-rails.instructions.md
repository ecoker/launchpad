---
name: Ruby on Rails
description: Convention over configuration — rapid full-stack web with incredible generators
applyTo: "**/*.{rb,erb,haml,slim}"
---

# Ruby on Rails

Rails is one of the best frameworks for AI-assisted development. Convention
over configuration means the AI doesn't have to make architectural decisions —
Rails already made them. The generator ecosystem is unmatched.

## Scaffold

```sh
rails new {{name}}
```

Use the Rails CLI and its generators extensively. Never hand-write migration
files, model boilerplate, or controller scaffolds.

### Key generators

```sh
# Full resource scaffold (model + migration + controller + views + routes)
rails generate scaffold Article title:string body:text published:boolean

# Model + migration only
rails generate model Comment body:text article:references

# Controller
rails generate controller Dashboard index

# Authentication (Rails 8+)
rails generate authentication
```

## Project structure

Rails defines the structure. Follow it:

```
app/
  controllers/        # Thin — delegate to models/services
  models/             # Business logic lives here
    concerns/         # Shared model behavior
  views/              # ERB templates
  services/           # Complex operations (use sparingly)
  jobs/               # Background processing
  mailers/
  channels/           # ActionCable for real-time
config/
  routes.rb           # RESTful routes
db/
  migrate/            # Versioned, reversible migrations
  schema.rb           # Auto-generated — never edit
test/                 # Minitest (default) or RSpec
```

## Rails conventions

### Models

- **Fat models, thin controllers.** Business logic belongs in models, not in
  controllers. Controllers orchestrate; models decide.
- **Use scopes for reusable queries.**
- **Validations at the model level.** Every user-facing write goes through
  model validations.
- **Use concerns judiciously.** Extract shared behavior into concerns, but
  don't create a concern for every method.

```ruby
# ✅ Clean model with scopes and validations
class Article < ApplicationRecord
  belongs_to :author, class_name: 'User'
  has_many :comments, dependent: :destroy

  validates :title, presence: true, length: { maximum: 200 }
  validates :body, presence: true

  scope :published, -> { where(published: true) }
  scope :recent, -> { order(created_at: :desc).limit(10) }

  def publish!
    update!(published: true, published_at: Time.current)
  end
end
```

### Controllers

- **RESTful actions only.** `index`, `show`, `new`, `create`, `edit`,
  `update`, `destroy`. If you need more, you probably need another controller.
- **Strong parameters always.** Never trust user input.
- **Before actions for shared setup.** `set_article` pattern for `show`,
  `edit`, `update`, `destroy`.

```ruby
# ✅ Thin, RESTful controller
class ArticlesController < ApplicationController
  before_action :set_article, only: %i[show edit update destroy]

  def index
    @articles = Article.published.recent
  end

  def create
    @article = current_user.articles.build(article_params)
    if @article.save
      redirect_to @article, notice: 'Article created.'
    else
      render :new, status: :unprocessable_entity
    end
  end

  private

  def set_article
    @article = Article.find(params[:id])
  end

  def article_params
    params.require(:article).permit(:title, :body, :published)
  end
end
```

### Services

Use service objects for multi-model operations:

```ruby
# ✅ Service for complex operations
class PublishArticleService
  def initialize(article, notify: true)
    @article = article
    @notify = notify
  end

  def call
    @article.publish!
    ArticleMailer.published(@article).deliver_later if @notify
    @article
  end
end
```

## Testing

- **Use Minitest** (Rails default) or RSpec — pick one, commit to it.
- **Model tests for business logic.**
- **System tests for critical user flows** (Capybara).
- **Use fixtures or factories** — not both.

## What to avoid

- Hand-writing migrations — use generators.
- Skipping validations with `update_column` or `save(validate: false)`.
- Business logic in controllers.
- N+1 queries — use `includes` and `strict_loading`.
- Callbacks for complex side effects — use service objects instead.
- API-only mode when you actually need views — Rails views are simple and
  AI-friendly.
