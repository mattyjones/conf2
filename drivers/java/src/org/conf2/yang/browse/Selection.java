package org.conf2.yang.browse;

import org.conf2.yang.Meta;
import org.conf2.yang.MetaCollection;
import org.conf2.yang.driver.DriverHandle;

/**
 *
 */
public class Selection {
    public MetaCollection meta;
    public Meta position;
    public BrowseEnter Enter;
    public BrowseRead Read;
    public BrowseIterate Iterate;
    public BrowseEdit Edit;
    public BrowseExit Exit;
    public BrowseChoose Choose;
    public boolean insideList;
    public boolean found;
}
