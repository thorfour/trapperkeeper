[![Docker Repository on Quay](https://quay.io/repository/thorfour/trapperkeeper/status "Docker Repository on Quay")](https://quay.io/repository/thorfour/trapperkeeper)

# trapperkeeper
slack integration that allows people to create blind submissions that can only be released after a a window of time has passed.

## Download

`docker pull quay.io/thorfour/trapperkeeper`

## Running

`docker run -d -p 80:80 -p 443:443 quay.io/thorfour/trapperkeeper /server -host <url> -email <support_email>`

### Using from customer slack integration

trapperkeeper supports the following commands:
  - window (check the current window expiration)
  - add [any text] (add a new submission)
  - release (release the current window if it has expired)
  - new [time] (create a new window with expiration i.e `new 1h` creates a 1 hour window)

## Building from source

`make docker` to create a standalone docker server

`make plugin` to creater a plugin (designed to work with github.com/thorfour/sillyputty)
