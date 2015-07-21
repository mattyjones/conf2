package org.conf2.meta;

/**
 *
 */
public class LinkedList {
    private Meta firstMeta;
    private Meta lastMeta;
    private MetaCollection parent;

    public LinkedList(MetaCollection parent) {
        this.parent = parent;
    }

    public Meta getFirstMeta() {
        return firstMeta;
    }

    public void addMeta(Meta m) {
        link(m);
    }

    public Meta getLastMeta() {
        return lastMeta;
    }

    public void link(Meta meta) {
        meta.setParent(parent);
        if (lastMeta != null) {
            lastMeta.setSibling(meta);
        }
        lastMeta = meta;
        if (firstMeta == null) {
            firstMeta = meta;
        }
    }

    public void setParent(MetaCollection parent) {
        this.parent = parent;
    }

    public MetaCollection getParent() {
        return parent;
    }
}
