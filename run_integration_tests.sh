#!/usr/bin/env bash

TEST_DIR=/tmp/combine-tests

createTestDir() {
  mkdir $TEST_DIR
  TEMPLATE_FILE="$TEST_DIR/tmpl.txt"
  TEMPLATE_CONTENTS="contents={{ .Read \"file.txt\" }}"
  echo $TEMPLATE_CONTENTS > $TEMPLATE_FILE

  FILE="$TEST_DIR/file.txt"
  FILE_CONTENTS='file'
  echo $FILE_CONTENTS > $FILE
  TEMPLATE_OUTPUT="contents=$FILE_CONTENTS"

  CSS_FILE="$TEST_DIR/css.txt"
  cat << EOF > $CSS_FILE

body {
  background: #ff0000;
}

EOF
  CSS_TEMPLATE_FILE="$TEST_DIR/csstmpl.txt"
  CSS_TEMPLATE_CONTENTS="{{ .Read \"css.txt\" }}"
  echo $CSS_TEMPLATE_CONTENTS > $CSS_TEMPLATE_FILE
  CSS_FILE_CONTENTS=$(cat "$TEST_DIR/css.txt")
  CSS_OUTPUT='body{background:red}'

  OUTPUT_FILE="$TEST_DIR/out.txt"
}

cleanup() {
  rm $FILE
  rm $CSS_FILE
  rm $CSS_TEMPLATE_FILE
  rm $TEMPLATE_FILE
  rmdir $TEST_DIR
}

fail() {
  echo ''
  echo 'FAIL'
  echo $1
}

pass() {
  echo -n '.'
}

doneWithTests() {
  echo ''
  echo 'All tests passed'
  exit 0
}

# Test that usage is printed if no flags is present
go run cmd/combine/main.go 2>/dev/null
STATUS=$?
if [[ $STATUS -ne 1 ]]
then
  fail 'Usage not printed with no flags'
  exit 1
fi
pass

# Test that base directory is required when reading from stdin
createTestDir
go run cmd/combine/main.go -o $OUTPUT_FILE <<< $TEMPLATE_CONTENTS 2>/dev/null
STATUS=$?
if [[ $STATUS -ne 1 ]]
then
  fail 'Base dir not required when reading from stdin'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
cleanup
pass

# Test that error happends when input is a broken template
createTestDir
go run cmd/combine/main.go -o $OUTPUT_FILE <<< "{{ file.Read \"file1.txt\" }" 2>/dev/null
STATUS=$?
if [[ $STATUS -ne 1 ]]
then
  fail 'Base dir not happen when input is a broken template'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
cleanup
pass

# Test that output is correct when reading from file and writing to file
createTestDir
go run cmd/combine/main.go -d $TEST_DIR -o $OUTPUT_FILE -i $TEMPLATE_FILE
OUTPUT=$(cat $OUTPUT_FILE)
if [[ $OUTPUT != $TEMPLATE_OUTPUT ]]
then
  fail '[Template] Output is not correct when reading from file and writing to file'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
rm $OUTPUT_FILE
cleanup
pass

# Test that output is correct when reading from stdin and writing to file
createTestDir
go run cmd/combine/main.go -d $TEST_DIR -o $OUTPUT_FILE <<< $TEMPLATE_CONTENTS
OUTPUT=$(cat $OUTPUT_FILE)
if [[ $OUTPUT != $TEMPLATE_OUTPUT ]]
then
  fail '[Template] Output is not correct when reading from stdin and writing to file'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
rm $OUTPUT_FILE
cleanup
pass

# Test that output is correct when reading from file and writing to stdout
createTestDir
OUTPUT=$(go run cmd/combine/main.go -i $TEMPLATE_FILE)
if [[ $OUTPUT != $TEMPLATE_OUTPUT ]]
then
  fail '[Template] Output is not correct when reading from file and writing to stdout'
  cleanup
  exit 1
fi
cleanup
pass

# Test that output is correct when reading from stdin and writing to stdout
createTestDir
OUTPUT=$(go run cmd/combine/main.go -d $TEST_DIR <<< $TEMPLATE_CONTENTS)
if [[ $OUTPUT != $TEMPLATE_OUTPUT ]]
then
  fail '[Template] Output is not correct when reading from stdin and writing to stdout'
  echo $OUTPUT
  exit 1
fi
cleanup
pass

# Test that output is correct when reading from stdin and writing to stdout and using a non defined type
createTestDir
OUTPUT=$(go run cmd/combine/main.go -d $TEST_DIR -t txt <<< $CSS_TEMPLATE_CONTENTS)
if [[ $OUTPUT != $CSS_FILE_CONTENTS ]]
then
  fail '[Template] Output is not correct when reading from stdin and writing to stdout and using a non defined type'
  echo $OUTPUT
  exit 1
fi
cleanup
pass

# Test yuiminifycombiner only if yuicompressor is installed
hash yuicompressor &> /dev/null
if [ $? -eq 1 ]; then
    echo 'Skipping tests for yuiminifycombiner as yuicompressor is not installed'
    doneWithTests
fi

# Test that output is correct when reading and minifying from stdin and writing to file
createTestDir
go run cmd/combine/main.go -d $TEST_DIR -o $OUTPUT_FILE -t css <<< $CSS_TEMPLATE_CONTENTS
OUTPUT=$(cat $OUTPUT_FILE)
if [[ $OUTPUT != $CSS_OUTPUT ]]
then
  fail '[Minify] Output is not correct when reading and minifying from stdin and writing to file'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
rm $OUTPUT_FILE
cleanup
pass

# Test that output is correct when reading from file and writing to file
createTestDir
go run cmd/combine/main.go -d $TEST_DIR -o $OUTPUT_FILE -i $CSS_TEMPLATE_FILE -t css
OUTPUT=$(cat $OUTPUT_FILE)
if [[ $OUTPUT != $CSS_OUTPUT ]]
then
  fail '[Minify] Output is not correct when reading from stdin and writing to file'
  rm $OUTPUT_FILE
  cleanup
  exit 1
fi
rm $OUTPUT_FILE
cleanup
pass

# Test that output is correct when reading from file and writing to stdout
createTestDir
OUTPUT=$(go run cmd/combine/main.go -i $CSS_TEMPLATE_FILE -t css)
if [[ $OUTPUT != $CSS_OUTPUT ]]
then
  fail '[Minify] Output is not correct when reading from file and writing to stdout'
  cleanup
  exit 1
fi
cleanup
pass

# Test that output is correct when reading from stdin and writing to stdout
createTestDir
OUTPUT=$(go run cmd/combine/main.go -d $TEST_DIR -t css <<< $CSS_TEMPLATE_CONTENTS)
if [[ $OUTPUT != $CSS_OUTPUT ]]
then
  fail '[Minify] Output is not correct when reading from stdin and writing to stdout'
  cleanup
  exit 1
fi
cleanup
pass

doneWithTests