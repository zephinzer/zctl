install:
	@printf -- '\033[32minstalling 'zctl'... \033[0m'
	@go install ./cmd/z \
		&& printf -- "\033[32mDONE!\n> add 'alias z=zctl' to your .bashrc or equivalent to use 'z' instead of 'zctl'\033[0m\n" \
		|| printf -- "\033[31mFAILED!\n> add 'failed to install zctl, see above logs for more information\033[0m\n"
