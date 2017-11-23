# rdb2om
rdb2om - Relational Database to Object Mapper

This Project is in an early development.

This Project has the goal to provide an object-like interface to relational databases. It includes eager loading, either via manual mapping or via automatic detection (maybe possible in some cases). To not run into endless recursion loops, it should contain a loop counter and a detection, which "nodes" has already been visited.

## Why not use an orm?

Some ORMs are very powerfull, but they exist only for some few programming languages (JRE and .NET languages), so not every programming language has access to a ORM with feature X. You could use a very powerfull ORM on another programming language than your project's language, but this means double class definitions, redundant code, slower performance and other downsides. Why not create an access layer to a database which lets you store your objects independend from an ORM? The downside is, you need additional mapping.

## Plan

- Webservice
- Cleanup, reorganize code
- Multiple Databases (plugins via dialects)
- Database Query via URL/POST
- In next development phase write/definition operations
- Manual mapping
- Create/Use a SQL Like Language in further stages.


## Similar Projects

- [Postgrest](https://postgrest.com/)
- [GORM Recursive fetcher](https://github.com/zonkiie/gorm_recursive_fetcher)

## Contributors needed

Please contact me if you are interested in contributing to this project.
