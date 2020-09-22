import sys
import json
import os
import logging
import pymysql
import boto3
import csv
from io import StringIO

logging.basicConfig()
logger = logging.getLogger() 
logger.setLevel(logging.INFO)


def get_db_password():
    db_pass = os.getenv('DB_PASS')

    if db_pass:
        return db_pass

    local = os.getenv('ENV', True)

    if local == True:
        return os.getenv('DB_PASS', 'root')

    ssm = boto3.client('ssm')
    pass_name = os.environ["SSM_DB_PASS_NAME"]

    try: 
        ssm_param = ssm.get_parameter(Name=pass_name, WithDecryption=True)
    except:
        logger.error("ERROR: Unexpected error: Could not retrieve password from ssm")
        sys.exit()

    return ssm_param['Parameter']['Value']

def handler(event, context):

    print(event)

    rds_host = os.environ["DB_HOST"]
    name = os.environ["DB_USER"]
    db_name = os.environ["DB_NAME"]
    password=get_db_password().strip()
    port = 3306

    try:
        conn = pymysql.connect(rds_host, user=name,
                           passwd=password, db=db_name, connect_timeout=5, charset='utf8mb4')
    except:
        logger.error("ERROR: Unexpected error: Could not connect to MySql instance.")
        sys.exit()

    with conn.cursor() as cur:
        cur.execute(event["query"])
        res = cur.fetchall()
        conn.commit()
        lastrowid = cur.lastrowid
        affected_rows = conn.affected_rows()

        # INSERT
        if lastrowid != 0 and lastrowid != None:
            return {'result': str(lastrowid)} 
        
        # update
        if lastrowid == 0 and affected_rows > 0:
            return {'result': 'OK'}
        
        # update
        if lastrowid == 0 and affected_rows == 0:
            return {'result': 'OK'}

        # none affected
        if affected_rows == 0:
            return {'result': 'None'}



        headers = [d[0] for d in cur.description]

        writer = csv.writer(sys.stdout)
        writer.writerow(headers)
        writer.writerows(res)

        data = StringIO()
        writer = csv.writer(data)
        writer.writerow(headers)
        writer.writerows(res)

        return {'result': str(data.getvalue())}
        
    
if __name__ == '__main__':
    handler({"query": sys.argv[1]}, {})
