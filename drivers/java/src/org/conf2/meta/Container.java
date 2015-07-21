package org.conf2.meta;

/**
 *
 */
public class Container extends CollectionBase implements Describable {
    private LinkedList groupings = new LinkedListCollection("groupings", this);
    private LinkedList typedefs = new LinkedListCollection("typedefs", this);
    private LinkedList choices = new LinkedListCollection("choices", this);

    public Container(String ident) {
        super(ident);
    }

    @Override
    public void addMeta(Meta m) {
        if (m instanceof Grouping) {
            groupings.addMeta(m);
        } else if (m instanceof Choice) {
            choices.addMeta(m);
        } else if (m instanceof Typedef) {
            typedefs.addMeta(m);
        } else {
            super.addMeta(m);
        }
    }
}
