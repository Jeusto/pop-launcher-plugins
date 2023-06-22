export function cleanSearchQuery(query: string, plugin_id: string): string {
  return query.trim().replace(new RegExp(`${plugin_id}(\\s|$)`), "");
}
