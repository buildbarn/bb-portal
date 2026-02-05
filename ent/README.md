# Ent models

This directory contains the Ent schema definitions used by the application.

- Schemas live in `ent/schema`.
- The final schema (with privacy policies) is produced by re-exporting those
  definitions from `ent/authschema`.

Why the two-step process

- Privacy policies import generated code. Defining policies in the same package
  as the schema creates a circular dependency during code generation making it
  hard to bootstrap the generated files.
- To avoid that, we generate files this way:
  1. Remove all generated files.
  2. Generate files from `ent/schema` without privacy policies.
  3. Generate files from `ent/authschema` to create the final files.

## Adding a new schema

To add a new schema:

1. Create a new file in `ent/schema` defining the schema without privacy
   policies.
2. Add an import and re-export the schema in `ent/authschema/schema.go`.
3. Add any necessary privacy policies in `ent/authschema/privacy.go`. Privacy
   policies are needed if:
    - The entity is directly exposed via the API such as `findBazelInvocations`
      and there is no other form of access control.
    - The entity can be accessed via a relation from a entity that can access
      entites across the instanceName boundary, such as accessing Invocations
      from another instanceName via a `User` entity. In this case,
      `BazelInvocation` needs a privacy policy.
4. Add tests for the privacy policies in `ent/authschema/privacy_test.go`.
5. Run `./tools/update_schema.sh` to regenerate the code with the new schema
   and policies.
