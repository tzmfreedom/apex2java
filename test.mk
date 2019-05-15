.SUFFIXES: .java .class
.java.class:
	@javac $*.java

.PHONY: run
run: classes
	@java Hoge fooo

.PHONY: build
build: classes


CLASSES = $(shell ls *.java com/freedom_man/*.java)

classes: $(CLASSES:.java=.class)

clean:
	$(RM) $(CLASSES:.java=.class)
