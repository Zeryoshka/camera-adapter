SERVICES_PATH='/etc/systemd/system'
# SERVICES_PATH=$(pwd)

path_to_file=$(pwd)

for tepmlate_name in $(ls services-tmpl | grep -E "*.tmpl"); do

    name=$(echo $tepmlate_name | sed 's/\.[^.]*$//')
    if [ ! -f "$SERVICES_PATH/$name" ]; then
        echo "Create service for \"$SERVICES_PATH/$name\""

        escaped_path=$(echo $path_to_file | sed 's/\//\\\//g')
        sed 's/{path-tmpl}/'$escaped_path'/g' \
            <$(pwd)/services-tmpl/$tepmlate_name \
            >$SERVICES_PATH/$name
    else
        echo "Service \"$SERVICES_PATH/$name\" already created"
    fi
    systemctl daemon-reload
    systemctl restart $name

done
