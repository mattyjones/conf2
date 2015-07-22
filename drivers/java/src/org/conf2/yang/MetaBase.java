package org.conf2.yang;

/**
 *
 */
public abstract class MetaBase implements Meta {
    private String ident;
    private String description;
    private MetaCollection parent;
    private Meta sibling;

    MetaBase(String ident) {
        this.ident = ident;
    }

    @Override
    public void setParent(MetaCollection parent) {
        this.parent = parent;
    }

    @Override
    public void setSibling(Meta sibling) {
        this.sibling = sibling;
    }

    @Override
    public Meta getSibling() {
        return sibling;
    }

    @Override
    public MetaCollection getParent() {
        return parent;
    }

    @Override
    public String getIdent() {
        return ident;
    }

    // subclass must declared "implements Describable" for this to be properly exposed
    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }
}
