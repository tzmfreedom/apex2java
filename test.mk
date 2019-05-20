.SUFFIXES: .java .class
.java.class:
	@javac $*.java

.PHONY: run
run: build
	@java -classpath . Hoge fooo

.PHONY: build
build: classes
	@javac -classpath . Hoge.java

CLASSES = $(shell ls com/freedom_man/system/*.java)

classes: $(CLASSES:.java=.class)

clean:
	$(RM) $(CLASSES:.java=.class)
