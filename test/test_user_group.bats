# vim: set ft=sh:
# -*- mode: sh -*-

@test "set user and group by name" {
  run ./go-jail -user "daemon" -group "daemon" -- id
	echo "output = ${output}"
	[[ "$output" = *"uid=1(daemon)"* ]]
	[[ "$output" = *"gid=1(daemon)"* ]]
}

@test "set user and group by ID" {
  run ./go-jail -user "#1" -group "#1" -- id
	echo "output = ${output}"
	[[ "$output" = *"uid=1(daemon)"* ]]
	[[ "$output" = *"gid=1(daemon)"* ]]
}
