package org.conf2.meta;

/**
 *
 */
public interface Meta extends Identifiable {
    void setParent(MetaCollection parent);
    void setSibling(Meta sibling);
    Meta getSibling();
    MetaCollection getParent();
}
