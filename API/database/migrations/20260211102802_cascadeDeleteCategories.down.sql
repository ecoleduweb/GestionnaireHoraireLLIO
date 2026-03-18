-- Migration Down
ALTER TABLE categories
  DROP FOREIGN KEY fk_categories_projects;

ALTER TABLE categories
  ADD CONSTRAINT fk_categories_projects
  FOREIGN KEY (project_id) REFERENCES projects(id);
