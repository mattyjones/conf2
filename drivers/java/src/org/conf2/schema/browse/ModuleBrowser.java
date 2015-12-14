package org.conf2.schema.browse;

import org.conf2.schema.*;
import org.conf2.schema.yang.YangModule;

import java.util.Iterator;

/**
 *
 */
public class ModuleBrowser implements Browser {
    private Module yang = YangModule.YANG;
    public Module module;

    public ModuleBrowser(Module module) {
        this.module = module;
    }

    public MetaCollection getSchema() {
        return yang;
    }

    @Override
    public Node getNode() {
        MyNode n = new MyNode();
        n.Enter = (Selection s, MetaCollection meta, boolean create) -> {
            if (create) {
                module = new Module("unknown");
            }
            if (module != null) {
                return enterModule(module);
            }
            return null;
        };
        n.Edit = (Selection s, MetaCollection meta, BrowseValue val) -> {
            BrowseUtil.setterMethod(meta, module, val);
        };
        return n;
    }

    Node enterModule(final Module module) {
        MyNode s = new MyNode();
        s.Read = (Selection sel, Meta meta) ->
                BrowseUtil.getterMethod(s.position, module);
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if (ident.equals("revision")) {
                if (create) {
                    module.revision = new Revision("unknown");
                }
                if (module.revision != null) {
                    return enterRevision(module.revision);
                }
            } else {
                return metaCollectionEnter(module, ident);
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
            BrowseUtil.setterMethod(meta, module, val);
        return s;
    }

    Node enterLeaf(Leaf leaf) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if ("type".equals(ident)) {
                if (create) {
                    leaf.setDataType(new DataType("unknown"));
                }
                if (leaf.getDataType() != null) {
                    return enterDataType(leaf.getDataType());
                }
                return null;
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, leaf, val);
        return s;
    }

    Node enterDataType(DataType type) {
        MyNode s = new MyNode();
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) -> {
            String ident = meta.getIdent();
            if ("ident".equals(ident)) {
                type.setIdent(val.str);
            } else {
                BrowseUtil.setterField(s.position, type, val);
            }
        };

        return s;
    }

    Node enterLeafList(LeafList leafList) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if ("type".equals(ident)) {
                if (create) {
                    leafList.setDataType(new DataType("unknown"));
                }
                if (leafList.getDataType() != null) {
                    return enterDataType(leafList.getDataType());
                }
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, leafList, val);
        return s;
    }

    Node enterUses(Uses uses) {
        MyNode s = new MyNode();
        // TODO s.Enter
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, uses, val);
        return s;
    }

    Node enterChoice(Choice choice) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if ("cases".equals(ident)) {
                return enterDefinitionCollection(choice);
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
            BrowseUtil.setterMethod(s.position, choice, val);
        return s;
    }

    Node enterDefinition(MetaCollection c, Meta m) {
        MyNode n = new MyNode();
        final Meta[] ptr = new Meta[] { m };
        n.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if (create) {
                switch (ident) {
                    case "leaf":
                        ptr[0] = new Leaf("unknown");
                        break;
                    case "leaf-list":
                        ptr[0] = new LeafList("unknown");
                        break;
                    case "uses":
                        ptr[0] = new Uses("unknown");
                        break;
                    case "choice":
                        ptr[0] = new Choice("unknown");
                        break;
                    case "action":
                    case "rpc":
                        ptr[0] = new Rpc("unknown");
                        break;
                    case "container":
                        ptr[0] = new Container("unknown");
                        break;
                    case "grouping":
                        ptr[0] = new Grouping("unknown");
                        break;
                    case "typedef":
                        ptr[0] = new Typedef("unknown");
                        break;
                    default:
                        throw new RuntimeException("unknown creation type " + ident);
                }
            }
            if (ptr[0] != null) {
                switch (ident) {
                    case "leaf":
                        return enterLeaf((Leaf)ptr[0]);
                    case "leaf-list":
                        return enterLeafList((LeafList)ptr[0]);
                    case "uses":
                        return enterUses((Uses)ptr[0]);
                    case "choice":
                        return enterChoice((Choice)ptr[0]);
                    case "action":
                    case "rpc":
                        return enterRpc((Rpc)ptr[0]);
                    case "container":
                        return enterContainer((Container) ptr[0]);
                    case "grouping":
                        return enterGrouping((Grouping)ptr[0]);
                    case "typedef":
                        return enterTypedef((Typedef)ptr[0]);
                    default:
                        throw new RuntimeException("unknown enter type " + ident);
                }
            }
            return null;
        };
        n.Event = (Selection sel, DataEvent e) -> {
            switch (e) {
                case NEW:
                    c.addMeta(ptr[0]);
                    break;
            }
        };
        n.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(meta, ptr[0], val);
        return n;
    }

    Node enterDefinitionCollection(MetaCollection c) {
        MyNode s = new MyNode();
        final Iterator<Meta>[] iPtr = new Iterator[1];
        s.Iterate = (Selection sel, BrowseValue[] key, boolean create, boolean isFirst) -> {
            if (create) {
                return enterDefinition(c, null);
            }
            Meta m = null;
            if (key.length > 0) {
                MetaUtil.findByIdent(c, key[0].str);
            } else {
                if (isFirst) {
                    iPtr[0] = new MetaCollectionIterator(c);
                }
                if (iPtr[0].hasNext()) {
                    m = iPtr[0].next();
                }
            }
            return null;
        };
        return s;
    }

    Node enterChoiceCase(ChoiceCase ccase) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            return metaCollectionEnter(ccase, meta.getIdent());
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, ccase, val);
        return s;
    }

    Node enterNotification(Notification notification) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if (ident.equals("groupings")) {
                return enterDefinitionCollection(notification.getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterDefinitionCollection(notification.getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitionCollection(notification);
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, notification, val);
        return s;
    }

    Node enterContainer(Container container) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            return metaCollectionEnter(container, meta.getIdent());
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
            BrowseUtil.setterMethod(s.position, container, val);
        return s;
    }

    Node metaCollectionEnter(MetaCollection c, String ident) {
        if (ident.equals("groupings")) {
            return enterDefinitionCollection(((HasGroupings)c).getGroupings());
        } else if (ident.equals("typedefs")) {
            return enterDefinitionCollection(((HasTypedefs)c).getTypedefs());
        } else if (ident.equals("definitions")) {
            return enterDefinitionCollection(c);
        } else if (ident.equals("rpcs") || ident.equals("actions")) {
            return enterDefinitionCollection(((HasActions)c).getRpcs());
        } else if (ident.equals("notifications")) {
            return enterDefinitionCollection(((HasNotifications)c).getNotifications());
        }
        return null;
    }

    Node enterGrouping(Grouping grouping) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            return metaCollectionEnter(grouping, meta.getIdent());
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(meta, grouping, val);
        return s;
    }

    Node enterTypedef(Typedef typedef) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if ("type".equals(ident)) {
                if (create) {
                    typedef.setDataType(new DataType("unknown"));
                }
                if (typedef.getDataType() != null) {
                    return enterDataType(typedef.getDataType());
                }
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(meta, typedef, val);
        return s;
    }

    Node enterList(MetaList list) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            return metaCollectionEnter(list, meta.getIdent());
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                    BrowseUtil.setterMethod(s.position, list, val);
        return s;
    }

    Node enterRpc(Rpc rpc) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            String ident = meta.getIdent();
            if ("input".equals(ident)) {
                if (create) {
                    rpc.input = new RpcInput();
                }
                if (rpc.input != null) {
                    return enterRpcBase(rpc.input);
                }
            } else if ("output".equals(ident)) {
                if (create) {
                    rpc.output = new RpcOutput();
                }
                if (rpc.output != null) {
                    return enterRpcBase(rpc.output);
                }
            }
            return null;
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, rpc, val);
        return s;
    }

    Node enterRpcBase(RpcBase rpcBase) {
        MyNode s = new MyNode();
        s.Enter = (Selection sel, MetaCollection meta, boolean create) -> {
            return metaCollectionEnter(rpcBase, meta.getIdent());
        };
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) ->
                BrowseUtil.setterMethod(s.position, rpcBase, val);
        return s;
    }

    Node enterRevision(Revision rev) {
        MyNode s = new MyNode();
        s.Read = (Selection sel, Meta meta) ->
                BrowseUtil.getterMethod(s.position, module);
        s.Edit = (Selection sel, MetaCollection meta, BrowseValue val) -> {
            String ident = s.position.getIdent();
            if (ident.equals("rev-date")) {
                rev.setIdent(val.str);
            } else {
                BrowseUtil.setterMethod(s.position, module, val);
            }
        };
        return s;
    }
}
