package org.conf2.encoding;

import static java.lang.String.format;

/**
 *
 */
public class Schema {
    private Schema parent;
    private Schema sibling;
    private Schema firstChild;
    private Schema lastChild;
    private int size;
    private String ident;

    public Schema(String ident) {
        this.ident = ident;
    }

    public Schema(Schema parent, String ident) {
        this.parent = parent;
        this.ident = ident;
    }

    public Schema getParent() {
        return this.parent;
    }

    public String getIdent() {
        return ident;
    }

    public Schema getFirstChild() {
        return firstChild;
    }

    public Schema getChildByIdent(String ident) {
        Schema candidate = firstChild;
        while (candidate != null) {
            if (candidate.ident.equals(ident)) {
                return candidate;
            }
            candidate = candidate.sibling;
        }
        return null;
    }

    public Schema requireChildByIdent(String element) {
        Schema node = getChildByIdent(element);
        if (node == null) {
            String msg = format("Could not find node named '%s' on '%s'", element, ident);
            throw new InvalidPathException(msg);
        }
        return node;
    }

    public Schema addChild(Schema child) {
        child.parent = this;
        if (lastChild == null) {
            firstChild = child;
        } else {
            lastChild.sibling = child;
        }
        lastChild = child;
        size++;
        return child;
    }

    public int size() {
        return size;
    }
}
