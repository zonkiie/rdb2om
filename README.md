# rdb2om
rdb2om - Relational Database to Object Mapper

This Project is in an early development.

This Project has the goal to provide an object-like interface to relational databases. It includes eager loading, either via manual mapping or via automatic detection. To not run into endless recursion loops, it should contain a loop counter and a detection, which "nodes" has already been visited.

## Plan

- Webservice
- Cleanup, reorganize code
- Multiple Databases (plugins via dialects)
- Database Query via URL/POST
- In next development phase write/definition operations
- Manual mapping


## Similar Projects

- [Postgrest](https://postgrest.com/)
- [GORM Recursive fetcher](https://github.com/zonkiie/gorm_recursive_fetcher)

