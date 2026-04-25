# trafficlightsim

This is a repository to see if we can make a trafficlight using "OpenGL" to have more realistich lightning effects as compared to basic drawing with pixels and lines.

## developer notes

Some notes on development in this chapter. The code is supposed to work when you have a recent version of golang installed. We are using a provided openGL library and vendored that into the repository.

### precommit usage

`pre-commit` has been setup for this project for basic linting rules. Run pre-commit as follows

```bash
pre-commit run -a
```

### Help from an agent

When not familiar with "OpenGL" as a programming paradigm we can make use of guidance from a agentic model. As an example exercise this project was using a cloud model via ollama.

- Setup ollama and run it (not explained here, you can find this elsewhere)
- Run an agent tool, for example:

```bash
ollama launch claude --model qwen3-coder:480b-cloud
```

```bash
ollama launch claude --model qwen3.5:cloud
```

At time of writing these notes, there a few free models that can still be used via ollama cloud, as long as you do not heavily use the free models it can help when unfamiliar with how to alter the created openGl code. If you have a machine with a decent GPU, you could better be off using a locally installed large language model instead.
