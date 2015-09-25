package org.conf2.schema.browse;

import org.conf2.schema.*;
import org.conf2.schema.driver.DriverError;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Collection;
import java.util.Iterator;

/**
 *
 */
public class BrowseUtil {

    public static final String accessorMethodNameFromMeta(String prefix, String ident) {
        StringBuilder sb = new StringBuilder(prefix.length() + ident.length());
        boolean upper = false;
        if (prefix.length() > 0) {
            upper = true;
            sb.append(prefix);
        }
        for (int i = 0; i < ident.length(); i++) {
            char c = ident.charAt(i);
            if (c == '-') {
                upper = true;
            } else {
                if (upper) {
                    sb.append(Character.toUpperCase(c));
                } else {
                    sb.append(c);
                }
                upper = false;
            }
        }
        return sb.toString();
    }

    public static final void setterMethod(Meta m, Object o, BrowseValue v) {
        String methodName = accessorMethodNameFromMeta("set", m.getIdent());
        try {
            Method method;
            switch (v.valType) {
                case INT32_LIST:
                    method = o.getClass().getMethod(methodName, int[].class);
                    method.invoke(o, new Object[] {v.int32list});
                    break;
                case STRING_LIST:
                    method = o.getClass().getMethod(methodName, String[].class);
                    method.invoke(o, new Object[] {v.strlist});
                    break;
                case BOOLEAN_LIST:
                    method = o.getClass().getMethod(methodName, boolean[].class);
                    method.invoke(o, new Object[] {v.boollist});
                    break;
                case INT32:
                    method = o.getClass().getMethod(methodName, int.class);
                    method.invoke(o, v.int32);
                    break;
                case STRING:
                    method = o.getClass().getMethod(methodName, String.class);
                    method.invoke(o, v.str);
                    break;
                case BOOLEAN:
                    method = o.getClass().getMethod(methodName, boolean.class);
                    method.invoke(o, v.bool);
                    break;
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Method %s not found on class %s", methodName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
    }

    public static final BrowseValue readField(Meta m, Object o) {
        String fieldName = accessorMethodNameFromMeta("", m.getIdent());
        BrowseValue v;
        try {
            DataType t = ((HasDataType) m).getDataType();
            v = new BrowseValue();
            ValueType valType = ((HasDataType) m).getDataType().valType;
            Field field = o.getClass().getField(fieldName);
            switch (valType) {
                case ENUMERATION:
                case INT32: {
                    v.int32 = field.getInt(o);
                    break;
                }
                case BOOLEAN:
                    v.bool = field.getBoolean(o);
                    break;
                case STRING:
                    v.str = field.get(o).toString();
                    break;
                case INT32_LIST:
                    v.int32list = (int[]) field.get(o);
                    if (v.int32list == null) {
                        return null;
                    }
                    break;
                case BOOLEAN_LIST:
                    v.boollist = (boolean[]) field.get(o);
                    if (v.boollist == null) {
                        return null;
                    }
                case STRING_LIST:
                    v.strlist = coerseStringList(field.get(o));
                    if (v.strlist == null) {
                        return null;
                    }
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Field %s not found on class %s", fieldName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
        return v;
    }

    public static final void setterField(Meta m, Object o, BrowseValue v) {
        String fieldName = accessorMethodNameFromMeta("", m.getIdent());
        try {
            Field field = o.getClass().getField(fieldName);
            switch (v.valType) {
                case INT32_LIST:
                    field.set(o, v.int32list);
                case INT32:
                    field.setInt(o, v.int32);
                    break;
                case STRING_LIST:
                    field.set(o, v.strlist);
                case STRING:
                    field.set(o, v.str);
                    break;
                case BOOLEAN_LIST:
                    field.set(o, v.boollist);
                case BOOLEAN:
                    field.setBoolean(o, v.bool);
                    break;
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Field %s not found on class %s", fieldName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
    }

    public static String[] coerseStringList(Object o) {
        if (o instanceof Collection) {
            Collection c = (Collection) o;
            String[] strlist = new String[c.size()];
            Iterator items = c.iterator();
            for (int i = 0; items.hasNext(); i++) {
                Object item = items.next();
                strlist[i] = (item != null ? item.toString() : null);
            }
            return strlist;
        }

        return (String[])o;
    }

    public static final BrowseValue getterMethod(Meta m, Object o) {
        String methodName = accessorMethodNameFromMeta("get", m.getIdent());
        BrowseValue v;
        try {
            Method method = o.getClass().getMethod(methodName);
            Object result = method.invoke(o);
            if (result == null) {
                return null;
            }
            v = new BrowseValue();
            ValueType valType = ((HasDataType) m).getDataType().valType;
            switch (valType) {
                case INT32_LIST:
                    v.int32list = (int[]) result;
                    break;
                case INT32:
                    v.int32 = (Integer) result;
                    break;
                case STRING_LIST:
                    v.strlist = coerseStringList(result);
                    break;
                case STRING:
                    v.str = result.toString();
                    break;
                case BOOLEAN_LIST:
                    v.boollist = (boolean[]) result;
                    break;
                case BOOLEAN:
                    v.bool = (Boolean) result;
                    break;
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            throw new MetaError("Method not found", e);
        }
        return v;
    }

    public static final String[] strArrayAppend(String[] strlist, String s) {
        if (strlist == null) {
            return new String[] { s };
        }
        String[] strlistNew = Arrays.copyOf(strlist, strlist.length + 1);
        strlistNew[strlist.length - 1] = s;
        return strlist;
    }

    public static final int[] intArrayAppend(int[] intlist, int n) {
        if (intlist == null) {
            return new int[] { n };
        }
        int[] intlistNew = Arrays.copyOf(intlist, intlist.length + 1);
        intlistNew[intlistNew.length - 1] = n;
        return intlistNew;
    }
}
