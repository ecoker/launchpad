defmodule TestAppWeb.ChatLiveTest do
  use TestAppWeb.ConnCase

  import Phoenix.LiveViewTest

  describe "mount" do
    test "renders the join form", %{conn: conn} do
      {:ok, view, html} = live(conn, "/")

      assert html =~ "test-app chat"
      assert html =~ "Pick a username"
      assert has_element?(view, "button", "Join")
    end
  end

  describe "join" do
    test "shows the message form after joining", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "alice"}) |> render_submit()

      assert has_element?(view, "button", "Send")
      assert has_element?(view, ~s(input[name="body"]))
    end

    test "broadcasts a join notification", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "alice"}) |> render_submit()

      assert render(view) =~ "alice joined the room"
    end

    test "shows the leaderboard", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "alice"}) |> render_submit()

      assert render(view) =~ "Talky Talk Leaderboard"
      assert render(view) =~ "alice"
    end

    test "rejects empty username", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "  "}) |> render_submit()

      assert render(view) =~ "Username can&#39;t be blank"
    end
  end

  describe "send message" do
    test "broadcasts and displays a message", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "alice"}) |> render_submit()
      view |> form("form", %{body: "hello!"}) |> render_submit()

      assert render(view) =~ "hello!"
      assert render(view) =~ "alice"
    end

    test "awards points for letters typed", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/")

      view |> form("form", %{username: "alice"}) |> render_submit()
      view |> form("form", %{body: "hi"}) |> render_submit()

      # "hi" = 2 letters
      html = render(view)
      assert html =~ "2"
    end
  end
end
