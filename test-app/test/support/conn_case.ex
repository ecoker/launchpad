defmodule TestAppWeb.ConnCase do
  @moduledoc "Conveniences for testing Phoenix endpoints."

  use ExUnit.CaseTemplate

  using do
    quote do
      @endpoint TestAppWeb.Endpoint

      use TestAppWeb, :verified_routes

      import Plug.Conn
      import Phoenix.ConnTest
    end
  end

  setup _tags do
    {:ok, conn: Phoenix.ConnTest.build_conn()}
  end
end
