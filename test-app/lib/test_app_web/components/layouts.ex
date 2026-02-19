defmodule TestAppWeb.Layouts do
  @moduledoc "Layout components for the application shell."

  use TestAppWeb, :html

  embed_templates "layouts/*"
end
