package org.conf2.yang;

/**
 *
 */
public class LinkedListCollection extends LinkedList implements MetaCollection {
    private String ident;
    private Meta sibling;

    public LinkedListCollection(String ident, MetaCollection parent) {
        super(parent);
        this.ident = ident;
    }

    @Override
    public void setSibling(Meta sibling) {
        this.sibling = sibling;
    }

    @Override
    public Meta getSibling() {
        return null;
    }

    @Override
    public String getIdent() {
        return ident;
    }
}
