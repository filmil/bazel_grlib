package package_name is
    function value return string;
end package;

package body package_name is
    function value return string is
    begin
        return string'("Hello world");
    end;
end package body;
