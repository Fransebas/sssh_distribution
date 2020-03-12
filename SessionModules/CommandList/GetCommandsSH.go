package CommandList

var GetCommands = `
case "$PATH" in
  (*[!:]:) PATH="$PATH:" ;;
esac

set -f; IFS=:
for dir in $PATH; do
  set +f
  [ -z "$dir" ] && dir="."
  for file in "$dir"/*; do
    if [ -x "$file" ] && ! [ -d "$file" ]; then
      # printf '%s = %s\n' "${file##*/}" "$file"
      printf '%s\n' "${file##*/}"
    fi
  done
done`
