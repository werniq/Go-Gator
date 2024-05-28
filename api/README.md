- Feature/Project Name: news-aggregator

# Summary
One paragraph explanation of the feature/project.

# APIs design

Explain APIs you want to introduce.
- How will they work?
- What will be the input arguments? What will be the output?
- Provide examples of models which will be used.
- How will they interact with existing APIs?

### Workflow

API for news aggregator will use already existing modules for parsing XML/HTML/JSON data. 
It will either make GET/POST requests, or extract data from files, 
in order to get a list of bytes, that will be decoded into news object.

### Arguments

Input arguments will be parameters for parsing. (e.g. keywords, sources, ts-from etc.)
Output will be list of the news, filtered (or not) by arguments

### Models

Parser, ParsingInstructions, ParsingParameters, News

### Interaction between APIs

It will interact with existing APIs while making requests to gather data. 

