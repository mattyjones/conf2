package org.conf2.encoding;

/**
 *
 */
public class Browser {
    private Schema schema;
    private Object data;

    public Browser(Schema schema, Object data) {
        this.schema = schema;
        this.data = data;
    }

    Object searchRecursive(Schema parent, Selector path, int i, Object data, SelectorHandler r) {
        String[] elements = path.getElements();
        if (elements.length == i) {
            return data;
        }
        String ident = elements[i];
        Schema n = parent.requireChildByIdent(ident);
        Object result = r.nextLevel(n, ident, data);
        return searchRecursive(n, path, i + 1, result, r);
    }

    public Object search(Selector path, SelectorHandler r) {
        return searchRecursive(schema, path, 0, data, r);
    }
}
