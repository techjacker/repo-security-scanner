env:
	@python3 -m venv env

deps-update:
	@pip install -r requirements-to-freeze.txt --upgrade
	@pip freeze > requirements.txt

deps:
	@pip install -r requirements.txt
	@pre-commit install

clean:
	@pip uninstall -yr requirements.txt
	@pip freeze > requirements.txt

autopep8:
	@autopep8 . --recursive --in-place --pep8-passes 2000 --verbose

autopep8-stats:
	@pep8 --quiet --statistics .

test:
	@pytest tests -vv

.PHONY: deps lint test* debug clean
