package org.conf2.schema;

/**
 *
 */
public abstract class CollectionBase extends MetaBase implements MetaCollection {
    private LinkedList children = new LinkedList(this);

    public CollectionBase(String ident) {
        super(ident);
    }

    @Override
    public Meta getFirstMeta() {
        return children.getFirstMeta();
    }

    @Override
    public void addMeta(Meta m) {
        children.addMeta(m);
    }
}
