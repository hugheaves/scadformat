#!/bin/bash

#test_ids=("bosl2")
test_ids=("internal")

declare -A test_input_dir test_input_repo_url test_input_repo_branch test_expected_dir test_expected_repo_url test_expected_repo_branch

test_input_dir[internal]=internal/formatter/testdata/valid
test_expected_dir[internal]=internal/formatter/testdata/expected

test_input_repo_url[bosl2]="https://github.com/hugheaves/BOSL2"
test_input_repo_branch[bosl2]="master"
test_expected_repo_url[bosl2]="https://github.com/hugheaves/BOSL2"
test_expected_repo_branch[bosl2]="scadformat"

function assert {
    exit_code=$1
    shift 1
    "$@"
    if [[ $? != "${exit_code}" ]]; then
        echo "FAILED: $*"
        exit 1
    fi

}

function test_setup {
    test_dir="${temp_dir}/${test_id}"
    input_dir="${test_dir}/input"
    expected_dir="${test_dir}/expected"
    formatted_dir="${test_dir}/formatted"

    mkdir -p "${test_dir}"
    mkdir -p "${formatted_dir}"

    if [[ -z  "${test_input_dir[${test_id}]}" ]]; then
      assert 0 git clone "${test_input_repo_url[${test_id}]}" "${input_dir}"
      assert 0 git -C "${input_dir}" checkout "${test_input_repo_branch[${test_id}]}"
    else
      assert 0 cp -R "${test_input_dir[${test_id}]}" "${input_dir}"
    fi

    if [[ -z  "${test_expected_dir[${test_id}]}" ]]; then

      assert 0 git clone "${test_expected_repo_url[${test_id}]}" "${expected_dir}"
      assert 0 git -C "${expected_dir}" checkout "${test_expected_repo_branch[${test_id}]}"
    else
      assert 0 cp -R "${test_expected_dir[${test_id}]}" "${expected_dir}"
    fi

    assert 0 find "${input_dir}" -type f -not -name "*.scad" -delete
    assert 0 find "${expected_dir}" -type f -not -name "*.scad" -delete
}

function test_formatting {
    assert 0 pushd "${input_dir}"

    while read -r filename; do
        echo "Testing formatting of ${filename}"
        dir_name=$(dirname "${filename}")
        mkdir -p "${formatted_dir}/${dir_name}"
        echo "Formatting ${filename} to ${formatted_dir}/${filename}"
        assert 0 scadformat <"${filename}" >"${formatted_dir}/${filename}"
        assert 0 diff "${formatted_dir}/${filename}" "${expected_dir}/${filename}"
    done < <(find . -name "*.scad")

    assert 0 popd
}

function test_openscad {
    assert 0 pushd "${input_dir}"

    while read -r filename; do
        echo "Testing rendering of ${filename}"

        dir_name=$(dirname "${filename}")


        # Generate STL from unformatted scad
        echo "Rendering ${input_dir}/${filename}"
        pushd "${input_dir}/${dir_name}"
        openscad -o "${temp_dir}/${test_id}/unformatted.stl" "${input_dir}/${filename}"
        unformatted_exit_code=$?
        popd

        # Generate STL from formatted scad
        echo "Rendering ${formatted_dir}/${filename}"
        pushd "${formatted_dir}/${dir_name}"
        openscad -o "${temp_dir}/${test_id}/formatted.stl" "${formatted_dir}/${filename}"
        formatted_exit_code=$?
        popd

        if [[ "${formatted_exit_code}" != "${unformatted_exit_code}" ]]; then
            echo "openscad test failed for ${filename}"
            exit 1
        fi

        if [[ "${formatted_exit_code}" == "0" ]]; then
          assert 0 diff "${temp_dir}/${test_id}/formatted.stl" "${temp_dir}/${test_id}/unformatted.stl"
        fi

    done < <(find . -name "*.scad")

    assert 0 popd
}

function run_test {
    test_id=$1
    test_setup
    test_formatting
    test_openscad
}

function main {
    if ! command -v openscad >/dev/null; then
        echo "openscad executable could not be found"
        exit 1
    fi

    temp_dir=$(mktemp -d)
    echo "Using directory ${temp_dir}"
    for test_id in "${test_ids[@]}"; do
        run_test "${test_id}"
    done
}

main
